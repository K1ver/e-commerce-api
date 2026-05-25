package postgres

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order domain.Order) error {
	const queryCreateOrder = `INSERT INTO orders(user_id, total_price, status)
		VALUES ($1, $2, $3)
		Returning id;`
	err := r.db.QueryRowxContext(ctx, queryCreateOrder, order.UserId, order.TotalPrice, order.Status).Scan(&order.Id)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	const queryAddOrderItem = `INSERT INTO order_items(order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4)`
	for _, orderItem := range order.Items {
		_, err = r.db.ExecContext(ctx, queryAddOrderItem, orderItem.OrderId, orderItem.ProductId, orderItem.Price, orderItem.Quantity)
		if err != nil {
			return fmt.Errorf("add order item: %w", err)
		}
	}
	return nil
}

func (r *orderRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	const queryGetAllOrders = `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE user_id = $1`
	var orders []domain.Order
	err := r.db.SelectContext(ctx, &orders, queryGetAllOrders, userID)
	if err != nil {
		return nil, fmt.Errorf("get all orders: %w", err)
	}

	var orderItems []domain.OrderItem
	const queryGetOrderItemsByOrderID = `SELECT id, order_id, product_id, price, quantity FROM order_items WHERE order_id = $1`
	for i, order := range orders {
		err = r.db.SelectContext(ctx, &orderItems, queryGetOrderItemsByOrderID)
		if err != nil {
			return nil, fmt.Errorf("get order items: %w", err)
		}
		order.Items = orderItems
		orders[i] = order
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error {
	const query = `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update status in order: %w", err)
	}
	return nil
}
