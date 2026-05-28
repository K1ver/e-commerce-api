package service

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type OrderService interface {
	Create(ctx context.Context, order domain.Order) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error
}

type orderService struct {
	orderRepository postgres.OrderRepository
}

func NewOrderService(orderRepository postgres.OrderRepository) OrderService {
	return &orderService{orderRepository: orderRepository}
}

func (service *orderService) Create(ctx context.Context, order domain.Order) error {
	err := service.orderRepository.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("create order %w", err)
	}

	return nil
}

func (service *orderService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	return service.orderRepository.GetAllByUserID(ctx, userID)
}

func (service *orderService) UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error {
	return service.orderRepository.UpdateStatusByID(ctx, status, id)
}
