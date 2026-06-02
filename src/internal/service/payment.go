package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/K1ver/e-commerce-api/internal/config"
	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type PaymentService interface {
	CreateForOrder(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Payment, string, error)
	Sync(ctx context.Context, userID uuid.UUID, role domain.Role, paymentID uuid.UUID) (domain.PaymentStatus, error)
	GetByOrderID(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Payment, error)
}

type paymentService struct {
	paymentRepository postgres.PaymentRepository
	orderRepository   postgres.OrderRepository
	paymentHandler    *yookassa.PaymentHandler
	cfg               config.YooKassaConfig
	validate          *validator.Validate
}

func NewPaymentService(
	paymentRepository postgres.PaymentRepository,
	orderRepository postgres.OrderRepository,
	paymentHandler *yookassa.PaymentHandler,
	cfg config.YooKassaConfig,
	validate *validator.Validate,
) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		paymentHandler:    paymentHandler,
		cfg:               cfg,
		validate:          validate,
	}
}

func (s *paymentService) CreateForOrder(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Payment, string, error) {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, "", err
	}
	if role != domain.RoleAdmin && order.UserId != userID {
		return nil, "", domain.ErrForbidden
	}
	if order.Status != domain.OrderStatusPending {
		return nil, "", fmt.Errorf("order is not payable")
	}

	payment := domain.Payment{
		OrderID: orderID,
		Amount:  order.TotalPrice,
		Status:  domain.PaymentStatusPending,
	}

	confirmationURL, err := s.createExternalPayment(ctx, &payment)
	if err != nil {
		return nil, "", err
	}
	if err := s.paymentRepository.Create(ctx, payment); err != nil {
		return nil, "", err
	}
	return &payment, confirmationURL, nil
}

func (s *paymentService) createExternalPayment(ctx context.Context, payment *domain.Payment) (string, error) {
	if s.paymentHandler == nil {
		payment.ID = uuid.New()
		return "http://localhost/mock-payment/" + payment.ID.String(), nil
	}
	yopayment, err := s.paymentHandler.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    strconv.Itoa(payment.Amount),
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: s.cfg.ReturnURL,
		},
		Description: "Order payment",
	})

	if err != nil {
		return "", fmt.Errorf("payment create failed: %w", err)
	}
	payment.ID = uuid.MustParse(yopayment.ID)
	if yopayment.Confirmation != nil {
		if url := yopayment.Confirmation; url != "" {
			return url, nil
		}
	}
	return s.cfg.ReturnURL, nil
}

func (s *paymentService) Sync(ctx context.Context, userID uuid.UUID, role domain.Role, paymentID uuid.UUID) (domain.PaymentStatus, error) {
	payment, err := s.paymentRepository.GetByID(ctx, paymentID)
	if err != nil {
		return "", err
	}
	order, err := s.orderRepository.GetByID(ctx, payment.OrderID)
	if err != nil {
		return "", err
	}
	if role != domain.RoleAdmin && order.UserId != userID {
		return "", domain.ErrForbidden
	}

	var status domain.PaymentStatus
	if s.paymentHandler == nil {
		status = domain.PaymentStatusSuccess
	} else {
		yopayment, err := s.paymentHandler.FindPayment(ctx, paymentID.String())
		if err != nil {
			return "", fmt.Errorf("payment find failed: %w", err)
		}
		switch yopayment.Status {
		case "succeeded":
			status = domain.PaymentStatusSuccess
		case "canceled":
			status = domain.PaymentStatusFailed
		default:
			status = domain.PaymentStatusPending
		}
	}

	if err := s.paymentRepository.UpdateStatus(ctx, paymentID, status); err != nil {
		return "", err
	}
	if status == domain.PaymentStatusSuccess {
		_ = s.orderRepository.UpdateStatusByID(ctx, domain.OrderStatusPaid, order.Id)
	}
	return status, nil
}

func (s *paymentService) GetByOrderID(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Payment, error) {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleAdmin && order.UserId != userID {
		return nil, domain.ErrForbidden
	}
	payment, err := s.paymentRepository.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
