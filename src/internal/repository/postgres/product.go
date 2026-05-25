package postgres

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) error
	Update(ctx context.Context, product domain.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) error {
	const query = `INSERT INTO products(name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, query, product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
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
	const query = `SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return product, fmt.Errorf("find product by id: %w", err)
	}
	return product, nil

}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	const query = `SELECT id, name, description, price, stock, created_at, updated_at FROM products`
	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	return products, nil
}

func (r *productRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}
