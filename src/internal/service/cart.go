package service

import (
	"context"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type CartService interface {
	GetMine(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	AddItem(ctx context.Context, userID, productID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, userID, productID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, userID, productID uuid.UUID) error
	Clear(ctx context.Context, userID uuid.UUID) error
}

type cartService struct {
	cartRepository    postgres.CartRepository
	productRepository postgres.ProductRepository
}

func NewCartService(cartRepository postgres.CartRepository, productRepository postgres.ProductRepository) CartService {
	return &cartService{cartRepository: cartRepository, productRepository: productRepository}
}

func (cs *cartService) GetMine(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	return cs.cartRepository.GetWithItemsByUserID(ctx, userID)
}

func (cs *cartService) AddItem(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	if _, err := cs.productRepository.FindByID(ctx, productID); err != nil {
		return err
	}
	cart, err := cs.cartRepository.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return cs.cartRepository.AddItem(ctx, cart.ID, productID, quantity)
}

func (cs *cartService) UpdateItem(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	cart, err := cs.cartRepository.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return cs.cartRepository.UpdateItem(ctx, cart.ID, productID, quantity)
}

func (cs *cartService) RemoveItem(ctx context.Context, userID, productID uuid.UUID) error {
	cart, err := cs.cartRepository.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return cs.cartRepository.RemoveItem(ctx, cart.ID, productID)
}

func (cs *cartService) Clear(ctx context.Context, userID uuid.UUID) error {
	cart, err := cs.cartRepository.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return cs.cartRepository.Clear(ctx, cart.ID)
}
