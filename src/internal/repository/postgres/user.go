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
	GetCredentialsByUsername(ctx context.Context, username string) (uuid.UUID, string, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	const queryCreateUser = `
          INSERT INTO users (
                             first_name,
                             last_name,
                             username,
                             email,
                             password_hash
          )
          VALUES ($1, $2, $3, $4, $5)
          RETURNING id, created_at, updated_at`
	err := u.db.QueryRowxContext(ctx, queryCreateUser, user.FirstName, user.LastName, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
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

func (u *userRepository) GetCredentialsByUsername(ctx context.Context, username string) (uuid.UUID, string, error) {
	const query = `SELECT id, password_hash FROM users WHERE username = $1`
	var cred struct {
		ID           uuid.UUID `db:"id"`
		PasswordHash string    `db:"password_hash"`
	}
	err := u.db.GetContext(ctx, &cred, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", domain.ErrUserNotFound
		}
		return uuid.Nil, "", fmt.Errorf("get credentials: %w", err)
	}
	return cred.ID, cred.PasswordHash, nil
}

func (u *userRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `SELECT id, first_name, last_name, username, email, created_at, updated_at
	FROM users
	WHERE id = $1`
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

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `SELECT id, first_name, last_name, username, email, created_at, updated_at
	FROM users
	WHERE email = $1`
	var user domain.User
	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = `SELECT id, first_name, last_name, username, email, created_at, updated_at
	FROM users
	WHERE username = $1`
	var user domain.User
	err := u.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) error {
	const query = `UPDATE users SET first_name = $1, last_name = $2, username = $3, email = $4, updated_at = now() WHERE id = $5`
	_, err := u.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`
	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
