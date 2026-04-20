package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // создан, но не оплачен
	OrderStatusPaid      OrderStatus = "paid"      // оплачен
	OrderStatusShipped   OrderStatus = "shipped"   // отправлен
	OrderStatusCompleted OrderStatus = "completed" // доставлен
	OrderStatusCanceled  OrderStatus = "canceled"  // отменён
)

type Order struct {
	Id         uuid.UUID   `json:"id" db:"id"`
	UserId     uuid.UUID   `json:"userId" db:"user_id"`
	TotalPrice float64     `json:"totalPrice" db:"total_price" validate:"required,min=0"`
	Status     OrderStatus `json:"status" db:"status" validate:"required"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at" validate:"required"`
}

type OrderItem struct {
	Id        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	OrderId   uuid.UUID `json:"orderId" db:"order_id" validate:"required,uuid"`
	ProductId uuid.UUID `json:"productId" db:"product_id" validate:"required,uuid"`
	Price     float64   `json:"price" db:"price" validate:"required,min=0"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required,min=0"`
}
