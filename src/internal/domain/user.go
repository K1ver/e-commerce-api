package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")

type User struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	FirstName string    `json:"firstName" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"lastName" db:"last_name" validate:"required,min=2,max=50"`
	Username  string    `json:"username" db:"username" validate:"required,min=3,max=30,alphanum"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"password" db:"password_hash" validate:"required,min=6,max=100"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
