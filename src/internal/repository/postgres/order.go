package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	CreateFromCart(ctx context.Context, userID uuid.UUID) (*domain.Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
	UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	const queryCreateOrder = `INSERT INTO orders(user_id, total_price, status)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = tx.QueryRowxContext(ctx, queryCreateOrder, order.UserId, order.TotalPrice, order.Status).
		Scan(&order.Id, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	const queryAddOrderItem = `INSERT INTO order_items(order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4)`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryAddOrderItem, order.Id, item.ProductId, item.Price, item.Quantity)
		if err != nil {
			return fmt.Errorf("add order item: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *orderRepository) CreateFromCart(ctx context.Context, userID uuid.UUID) (*domain.Order, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var cartID uuid.UUID
	err = tx.GetContext(ctx, &cartID, `SELECT id FROM carts WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("get cart: %w", err)
	}

	var items []domain.CartItem
	err = tx.SelectContext(ctx, &items, `SELECT id, cart_id, product_id, quantity FROM cart_items WHERE cart_id = $1`, cartID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var total int64
	orderItems := make([]domain.OrderItem, 0, len(items))
	for _, item := range items {
		var product domain.Product
		err = tx.GetContext(ctx, &product, `SELECT id, price, stock FROM products WHERE id = $1 FOR UPDATE`, item.ProductId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrProductNotFound
			}
			return nil, fmt.Errorf("get product: %w", err)
		}
		if product.Stock < item.Quantity {
			return nil, domain.ErrInsufficientStock
		}
		_, err = tx.ExecContext(ctx, `UPDATE products SET stock = stock - $1, updated_at = now() WHERE id = $2`, item.Quantity, product.ID)
		if err != nil {
			return nil, fmt.Errorf("update stock: %w", err)
		}
		lineTotal := product.Price * int64(item.Quantity)
		total += lineTotal
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: product.ID,
			Price:     int(product.Price),
			Quantity:  item.Quantity,
		})
	}

	order := &domain.Order{
		UserId:     userID,
		TotalPrice: int(total),
		Status:     domain.OrderStatusPending,
		Items:      orderItems,
	}
	err = tx.QueryRowxContext(ctx,
		`INSERT INTO orders(user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		order.UserId, order.TotalPrice, order.Status,
	).Scan(&order.Id, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	for _, item := range orderItems {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items(order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4)`,
			order.Id, item.ProductId, item.Price, item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("add order item: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID)
	if err != nil {
		return nil, fmt.Errorf("clear cart: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return order, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	const query = `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE id = $1`
	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}
	var items []domain.OrderItem
	err = r.db.SelectContext(ctx, &items, `SELECT id, order_id, product_id, price, quantity FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	order.Items = items
	return &order, nil
}

func (r *orderRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	const query = `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	var orders []domain.Order
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all orders: %w", err)
	}
	for i, order := range orders {
		var orderItems []domain.OrderItem
		err = r.db.SelectContext(ctx, &orderItems, `SELECT id, order_id, product_id, price, quantity FROM order_items WHERE order_id = $1`, order.Id)
		if err != nil {
			return nil, fmt.Errorf("get order items: %w", err)
		}
		order.Items = orderItems
		orders[i] = order
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatusByID(ctx context.Context, status domain.OrderStatus, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE orders SET status = $1, updated_at = now() WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("update status in order: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}
