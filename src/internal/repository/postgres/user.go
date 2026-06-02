package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetCredentialsByUsername(ctx context.Context, username string) (uuid.UUID, string, domain.Role, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (first_name, last_name, username, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	err := u.db.QueryRowxContext(ctx, query,
		user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return domain.ErrEmailAlreadyExists
		}
		if strings.Contains(err.Error(), "users_username_key") {
			return domain.ErrUsernameAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (u *userRepository) GetCredentialsByUsername(ctx context.Context, username string) (uuid.UUID, string, domain.Role, error) {
	const query = `SELECT id, password_hash, role FROM users WHERE username = $1`
	var cred struct {
		ID           uuid.UUID   `db:"id"`
		PasswordHash string      `db:"password_hash"`
		Role         domain.Role `db:"role"`
	}
	err := u.db.GetContext(ctx, &cred, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", "", domain.ErrUserNotFound
		}
		return uuid.Nil, "", "", fmt.Errorf("get credentials: %w", err)
	}
	return cred.ID, cred.PasswordHash, cred.Role, nil
}

func (u *userRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `SELECT id, first_name, last_name, username, email, role, created_at, updated_at
		FROM users WHERE id = $1`
	var user domain.User
	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (u *userRepository) List(ctx context.Context) ([]domain.User, error) {
	const query = `SELECT id, first_name, last_name, username, email, role, created_at, updated_at FROM users ORDER BY created_at DESC`
	var users []domain.User
	if err := u.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (u *userRepository) UpdateRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	const query = `UPDATE users SET role = $1, updated_at = now() WHERE id = $2`
	res, err := u.db.ExecContext(ctx, query, role, id)
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`
	res, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
