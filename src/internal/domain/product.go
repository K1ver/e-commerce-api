package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id" db:"id" validate:"omitempty,uuid"`
	SellerID    *uuid.UUID `json:"sellerId,omitempty" db:"seller_id"`
	Name        string     `json:"name" db:"name" validate:"required,min=2,max=50"`
	Description string     `json:"description" db:"description" validate:"required,min=2,max=2000"`
	Price       int64      `json:"price" db:"price" validate:"required,gt=0"`
	Stock       int        `json:"stock" db:"stock" validate:"gte=0"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}
