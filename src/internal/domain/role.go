package domain

import "errors"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleSeller Role = "seller"
	RoleBuyer  Role = "buyer"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleSeller, RoleBuyer:
		return true
	default:
		return false
	}
}

var ErrForbidden = errors.New("forbidden")
var ErrProductNotFound = errors.New("product not found")
var ErrCartNotFound = errors.New("cart not found")
var ErrOrderNotFound = errors.New("order not found")
var ErrInsufficientStock = errors.New("insufficient stock")
var ErrInvalidRole = errors.New("invalid role")
