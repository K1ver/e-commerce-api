package service

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type OrderService interface {
	Checkout(ctx context.Context, userID uuid.UUID) (*domain.Order, error)
	GetByID(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Order, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error
}

type orderService struct {
	orderRepository postgres.OrderRepository
}

func NewOrderService(orderRepository postgres.OrderRepository) OrderService {
	return &orderService{orderRepository: orderRepository}
}

func (s *orderService) Checkout(ctx context.Context, userID uuid.UUID) (*domain.Order, error) {
	return s.orderRepository.CreateFromCart(ctx, userID)
}

func (s *orderService) GetByID(ctx context.Context, userID uuid.UUID, role domain.Role, orderID uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleAdmin && order.UserId != userID {
		return nil, domain.ErrForbidden
	}
	return order, nil
}

func (s *orderService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	return s.orderRepository.GetAllByUserID(ctx, userID)
}

func (s *orderService) UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error {
	switch status {
	case domain.OrderStatusPending, domain.OrderStatusPaid, domain.OrderStatusShipped,
		domain.OrderStatusCompleted, domain.OrderStatusCanceled:
	default:
		return fmt.Errorf("invalid order status")
	}
	return s.orderRepository.UpdateStatusByID(ctx, status, id)
}
