package service

import (
	"context"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, product domain.Product) error
	Update(ctx context.Context, product domain.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	productRepository postgres.ProductRepository
}

func NewProductService(productRepository postgres.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (ps *productService) Create(ctx context.Context, product domain.Product) error {
	return ps.productRepository.Create(ctx, product)
}

func (ps *productService) Update(ctx context.Context, product domain.Product) error {
	return ps.productRepository.Update(ctx, product)
}

func (ps *productService) FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return ps.productRepository.FindByID(ctx, id)
}

func (ps *productService) FindAll(ctx context.Context) ([]domain.Product, error) {
	return ps.productRepository.FindAll(ctx)
}

func (ps *productService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return ps.productRepository.DeleteByID(ctx, id)
}
