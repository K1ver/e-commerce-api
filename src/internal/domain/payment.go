package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending" // создан и ждёт оплаты
	PaymentStatusSuccess PaymentStatus = "success" // оплачен
	PaymentStatusFailed  PaymentStatus = "failed"  // не оплачен
)

type Payment struct {
	ID        uuid.UUID   `json:"id" db:"id" validate:"required,uuid"`
	Amount    float64     `json:"amount" db:"amount" validate:"required,min=0"`
	Status    OrderStatus `json:"status" db:"status" validate:"required"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at" validate:"required"`
}
