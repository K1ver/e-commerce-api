package postgres

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment domain.Payment) error
	Update(ctx context.Context, payment domain.Payment) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (p *paymentRepository) Create(ctx context.Context, payment domain.Payment) error {
	const query = `INSERT INTO payments(order_id, amount, status) VALUES ($1, $2, $3)`
	_, err := p.db.ExecContext(ctx, query, payment.OrderID, payment.Amount, payment.Status)
	if err != nil {
		return fmt.Errorf("create payment: %w", err)
	}
	return nil
}

func (p *paymentRepository) Update(ctx context.Context, payment domain.Payment) error {
	const query = `UPDATE payments SET status = $2, updated_at = now() WHERE order_id = $1`
	_, err := p.db.ExecContext(ctx, query, payment.OrderID, payment.Amount)
	if err != nil {
		return fmt.Errorf("update payment: %w", err)
	}
	return nil
}

func (p *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	const query = `SELECT order_id, amount, status, created_at, updated_at FROM payments WHERE order_id = $1`
	var payment *domain.Payment
	err := p.db.GetContext(ctx, payment, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get all payments: %w", err)
	}
	return payment, nil
}
