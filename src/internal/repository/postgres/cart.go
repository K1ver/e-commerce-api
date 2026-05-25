package postgres

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	Create(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	GetWithItems(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error)
	AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	Clear(ctx context.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	const query = `INSERT INTO carts(user_id)
		VALUES($1)
		RETURNING id, created_at, updated_at`
	var cart domain.Cart
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create cart: %w", err)
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	const query = `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart by user UserID: %w", err)
	}
	return &cart, nil
}

func (r *cartRepository) GetWithItems(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error) {
	const query = `SELECT id, user_id, created_at, updated_at
		FROM carts
		WHERE id = $1;`
	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("get cart by cart ID: %w", err)
	}
	const query2 = `SELECT id, cart_id, product_id, quantity
		FROM cart_items
		WHERE cart_id = $1;`
	var cartItems []domain.CartItem
	err = r.db.SelectContext(ctx, &cartItems, query2, cartID)
	if err != nil {
		return nil, fmt.Errorf("get cart items by cart ID: %w", err)
	}
	cart.Items = cartItems
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	const query = `INSERT INTO cart_items(cart_id, product_id, quantity)VALUES($1, $2, $3);`
	_, err := r.db.ExecContext(ctx, query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("add cart item: %w", err)
	}
	return nil
}

func (r *cartRepository) UpdateItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	const query = `UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3;`
	_, err := r.db.ExecContext(ctx, query, quantity, cartID, productID)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	const query = `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;`
	_, err := r.db.ExecContext(ctx, query, cartID, productID)
	if err != nil {
		return fmt.Errorf("remove item: %w", err)
	}
	return nil
}

func (r *cartRepository) Clear(ctx context.Context, cartID uuid.UUID) error {
	const query = `DELETE FROM cart_items WHERE cart_id = $1;`
	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}
	return nil
}
