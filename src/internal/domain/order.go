package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // Создан, но не оплачен
	OrderStatusPaid      OrderStatus = "paid"      // Оплачен
	OrderStatusShipped   OrderStatus = "shipped"   // Отправлен
	OrderStatusCompleted OrderStatus = "completed" // Доставлен
	OrderStatusCanceled  OrderStatus = "canceled"  // Отменён
)

type Order struct {
	Id         uuid.UUID   `json:"id" db:"id"`
	UserId     uuid.UUID   `json:"userId" db:"user_id"`
	TotalPrice int         `json:"totalPrice" db:"total_price" validate:"required,min=1"`
	Status     OrderStatus `json:"status" db:"status" validate:"required"`
	Items      []OrderItem `json:"items" db:"-"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time   `json:"updatedAt" db:"updated_at"`
}

type OrderItem struct {
	Id        uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	OrderId   uuid.UUID `json:"orderId" db:"order_id" validate:"required,uuid"`
	ProductId uuid.UUID `json:"productId" db:"product_id" validate:"required,uuid"`
	Price     int       `json:"price" db:"price" validate:"required,min=1"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required,min=1"`
}
