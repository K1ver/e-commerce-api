package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment domain.Payment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (p *paymentRepository) Create(ctx context.Context, payment domain.Payment) error {
	const query = `INSERT INTO payments(id, order_id, amount, status) VALUES ($1, $2, $3, $4)`
	_, err := p.db.ExecContext(ctx, query, payment.ID, payment.OrderID, payment.Amount, payment.Status)
	if err != nil {
		return fmt.Errorf("create payment: %w", err)
	}
	return nil
}

func (p *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	const query = `UPDATE payments SET status = $2, updated_at = now() WHERE id = $1`
	res, err := p.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("update payment: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

func (p *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	const query = `SELECT id, order_id, amount, status, created_at, updated_at FROM payments WHERE id = $1`
	var payment domain.Payment
	err := p.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("get payment: %w", err)
	}
	return &payment, nil
}

func (p *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	const query = `SELECT id, order_id, amount, status, created_at, updated_at FROM payments WHERE order_id = $1`
	var payment domain.Payment
	err := p.db.GetContext(ctx, &payment, query, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("get payment: %w", err)
	}
	return &payment, nil
}
