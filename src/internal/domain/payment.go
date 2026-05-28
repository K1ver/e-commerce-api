package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"   // создан и ждёт оплаты
	PaymentStatusSuccess PaymentStatus = "succeeded" // оплачен
	PaymentStatusFailed  PaymentStatus = "canceled"  // не оплачен
)

type Payment struct {
	ID        uuid.UUID     `json:"id" db:"id" validate:"required,uuid"`
	OrderID   uuid.UUID     `json:"orderId" db:"order_id" validate:"required,uuid"`
	Amount    int           `json:"amount" db:"amount" validate:"required,min=1"`
	Status    PaymentStatus `json:"status" db:"status" validate:"required"`
	CreatedAt time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time     `json:"updatedAt" db:"updated_at"`
}
