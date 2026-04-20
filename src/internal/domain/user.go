package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	FirstName string    `json:"firstName" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"lastName" db:"last_name" validate:"required,min=2,max=50"`
	Username  string    `json:"username" db:"username" validate:"required,min=3,max=30,alphanum"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"password" db:"password" validate:"required,min=6,max=100"`
	CreatedAt time.Time `json:"createdAt" db:"created_at" validate:"required"`
}
