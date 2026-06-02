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

type CartRepository interface {
	GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	GetWithItemsByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
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

func (r *cartRepository) GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	cart, err := r.getByUserID(ctx, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, domain.ErrCartNotFound) {
		return nil, err
	}
	const query = `INSERT INTO carts(user_id) VALUES($1) RETURNING id, user_id, created_at, updated_at`
	cart = &domain.Cart{}
	err = r.db.QueryRowxContext(ctx, query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create cart: %w", err)
	}
	return cart, nil
}

func (r *cartRepository) getByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	const query = `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("get cart: %w", err)
	}
	return &cart, nil
}

func (r *cartRepository) GetWithItemsByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	cart, err := r.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	const query = `SELECT id, cart_id, product_id, quantity FROM cart_items WHERE cart_id = $1`
	var items []domain.CartItem
	if err := r.db.SelectContext(ctx, &items, query, cart.ID); err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	cart.Items = items
	return cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	const query = `
		INSERT INTO cart_items(cart_id, product_id, quantity) VALUES($1, $2, $3)
		ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`
	_, err := r.db.ExecContext(ctx, query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("add cart item: %w", err)
	}
	return nil
}

func (r *cartRepository) UpdateItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3`, quantity, cartID, productID)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`, cartID, productID)
	if err != nil {
		return fmt.Errorf("remove item: %w", err)
	}
	return nil
}

func (r *cartRepository) Clear(ctx context.Context, cartID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID)
	if err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}
	return nil
}
