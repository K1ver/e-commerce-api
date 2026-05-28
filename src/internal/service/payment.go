package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type PaymentService interface {
	Create(ctx context.Context, payment domain.Payment) (string, error)
	Update(ctx context.Context, payment domain.Payment) (domain.PaymentStatus, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
}

type paymentService struct {
	paymentRepository postgres.PaymentRepository
	paymentHandler    *yookassa.PaymentHandler
}

func NewPaymentService(paymentRepository postgres.PaymentRepository, paymentHandler *yookassa.PaymentHandler) PaymentService {
	return &paymentService{paymentRepository: paymentRepository, paymentHandler: paymentHandler}
}

func (s *paymentService) Create(ctx context.Context, payment domain.Payment) (string, error) {
	yopayment, err := s.paymentHandler.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    strconv.Itoa(payment.Amount),
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: "localhost/orders",
		},
		Description: "Test payment",
	})

	if err != nil {
		return "", fmt.Errorf("payment create failed: %w", err)
	}

	payment.ID = uuid.MustParse(yopayment.ID)

	err = s.paymentRepository.Create(ctx, payment)
	if err != nil {
		_, _ = s.paymentHandler.CancelPayment(ctx, yopayment.ID)
		return "", fmt.Errorf("payment create in db failed: %w", err)
	}
	return yopayment.ID, nil
}

func (s *paymentService) Update(ctx context.Context, payment domain.Payment) (domain.PaymentStatus, error) {
	yopayment, _ := s.paymentHandler.FindPayment(ctx, payment.ID.String())
	if yopayment.Status == "succeeded" {
		payment.Status = domain.PaymentStatusSuccess
	} else if yopayment.Status == "canceled" {
		payment.Status = domain.PaymentStatusFailed
	} else {
		payment.Status = domain.PaymentStatusPending
	}

	err := s.paymentRepository.Update(ctx, payment)
	return payment.Status, err
}

func (s *paymentService) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	return s.paymentRepository.GetByOrderID(ctx, orderID)
}
