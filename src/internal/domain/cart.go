package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	UserID    uuid.UUID `json:"userId" db:"user_id" validate:"required,uuid"`
	CreatedAt time.Time `json:"createdAt" db:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at" validate:"required"`
}
type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	CartId    uuid.UUID `json:"cartId" db:"cart_id" validate:"required,uuid"`
	ProductId uuid.UUID `json:"productId" db:"product_id" validate:"required,uuid"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required,min=0"`
}
