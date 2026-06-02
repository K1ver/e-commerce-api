package handler

import (
	"time"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
)

type userResponse struct {
	ID        uuid.UUID   `json:"id"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
	}
}

type productResponse struct {
	ID          uuid.UUID  `json:"id"`
	SellerID    *uuid.UUID `json:"sellerId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int64      `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func toProductResponse(p domain.Product) productResponse {
	return productResponse{
		ID:          p.ID,
		SellerID:    p.SellerID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toProductsResponse(items []domain.Product) []productResponse {
	out := make([]productResponse, len(items))
	for i, p := range items {
		out[i] = toProductResponse(p)
	}
	return out
}
