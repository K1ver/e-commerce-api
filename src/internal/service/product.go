package service

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/go-playground/validator/v10"
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
	validate          *validator.Validate
}

func NewProductService(productRepository postgres.ProductRepository, validate *validator.Validate) ProductService {
	return &productService{productRepository: productRepository, validate: validate}
}

func (ps *productService) Create(ctx context.Context, product domain.Product) error {
	err := ps.validate.StructCtx(ctx, product)
	if err != nil {
		return fmt.Errorf("validate err: %w", err)
	}

	return ps.productRepository.Create(ctx, product)
}

func (ps *productService) Update(ctx context.Context, product domain.Product) error {
	err := ps.validate.StructCtx(ctx, product)
	if err != nil {
		return fmt.Errorf("validate err: %w", err)
	}

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
