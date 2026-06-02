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
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, actorID uuid.UUID, actorRole domain.Role, product domain.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindBySeller(ctx context.Context, sellerID uuid.UUID) ([]domain.Product, error)
	Delete(ctx context.Context, actorID uuid.UUID, actorRole domain.Role, id uuid.UUID) error
}

type productService struct {
	productRepository postgres.ProductRepository
	validate          *validator.Validate
}

func NewProductService(productRepository postgres.ProductRepository, validate *validator.Validate) ProductService {
	return &productService{productRepository: productRepository, validate: validate}
}

func (ps *productService) Create(ctx context.Context, product *domain.Product) error {
	if err := ps.validate.StructCtx(ctx, product); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return ps.productRepository.Create(ctx, product)
}

func (ps *productService) Update(ctx context.Context, actorID uuid.UUID, actorRole domain.Role, product domain.Product) error {
	if err := ps.validate.StructCtx(ctx, product); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	existing, err := ps.productRepository.FindByID(ctx, product.ID)
	if err != nil {
		return err
	}
	if err := ps.canManageProduct(actorID, actorRole, existing); err != nil {
		return err
	}
	return ps.productRepository.Update(ctx, product)
}

func (ps *productService) FindByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return ps.productRepository.FindByID(ctx, id)
}

func (ps *productService) FindAll(ctx context.Context) ([]domain.Product, error) {
	return ps.productRepository.FindAll(ctx)
}

func (ps *productService) FindBySeller(ctx context.Context, sellerID uuid.UUID) ([]domain.Product, error) {
	return ps.productRepository.FindBySellerID(ctx, sellerID)
}

func (ps *productService) Delete(ctx context.Context, actorID uuid.UUID, actorRole domain.Role, id uuid.UUID) error {
	existing, err := ps.productRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := ps.canManageProduct(actorID, actorRole, existing); err != nil {
		return err
	}
	return ps.productRepository.DeleteByID(ctx, id)
}

func (ps *productService) canManageProduct(actorID uuid.UUID, actorRole domain.Role, product domain.Product) error {
	if actorRole == domain.RoleAdmin {
		return nil
	}
	if actorRole == domain.RoleSeller && product.SellerID != nil && *product.SellerID == actorID {
		return nil
	}
	return domain.ErrForbidden
}
