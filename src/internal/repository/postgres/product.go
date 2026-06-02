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

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product domain.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindBySellerID(ctx context.Context, sellerID uuid.UUID) ([]domain.Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	const query = `INSERT INTO products(name, description, price, stock, seller_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query,
		product.Name, product.Description, product.Price, product.Stock, product.SellerID,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) Update(ctx context.Context, product domain.Product) error {
	const query = `UPDATE products SET name = $2, description = $3, price = $4, stock = $5, updated_at = now() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	const query = `SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, domain.ErrProductNotFound
		}
		return product, fmt.Errorf("find product by id: %w", err)
	}
	return product, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	const query = `SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC`
	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	return products, nil
}

func (r *productRepository) FindBySellerID(ctx context.Context, sellerID uuid.UUID) ([]domain.Product, error) {
	const query = `SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products WHERE seller_id = $1 ORDER BY created_at DESC`
	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("find products by seller: %w", err)
	}
	return products, nil
}

func (r *productRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
