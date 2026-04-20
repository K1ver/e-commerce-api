package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=50"`
	Description string    `json:"description" db:"description" validate:"required,min=2,max=200"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0"`
	Stock       int       `json:"stock" db:"stock" validate:"required,min=0"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at" validate:"required"`
}
