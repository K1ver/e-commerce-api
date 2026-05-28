package service

import (
	"context"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type CartService interface {
	Create(ctx context.Context, userID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	GetWithItems(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error)
	AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	Clear(ctx context.Context, cartID uuid.UUID) error
}

type cartService struct {
	cartRepository postgres.CartRepository
}

func NewCartService(cartRepository postgres.CartRepository) CartService {
	return &cartService{cartRepository: cartRepository}
}

func (cs *cartService) Create(ctx context.Context, userID uuid.UUID) error {
	return cs.cartRepository.Create(ctx, userID)
}

func (cs *cartService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	return cs.cartRepository.GetByUserID(ctx, userID)
}

func (cs *cartService) GetWithItems(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error) {
	return cs.cartRepository.GetWithItems(ctx, cartID)
}

func (cs *cartService) AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	return cs.cartRepository.AddItem(ctx, cartID, productID, quantity)
}

func (cs *cartService) UpdateItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) error {
	return cs.cartRepository.UpdateItem(ctx, cartID, productID, quantity)
}

func (cs *cartService) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	return cs.cartRepository.RemoveItem(ctx, cartID, productID)
}

func (cs *cartService) Clear(ctx context.Context, cartID uuid.UUID) error {
	return cs.cartRepository.Clear(ctx, cartID)
}
