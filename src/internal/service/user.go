package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	SignIn(ctx context.Context, username, password string) (*domain.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepository postgres.UserRepository
	validate       *validator.Validate
}

func NewUserService(userRepository postgres.UserRepository, validate *validator.Validate) UserService {
	return &userService{userRepository: userRepository, validate: validate}
}

func (us *userService) Create(ctx context.Context, user *domain.User) error {
	if user.Role == "" {
		user.Role = domain.RoleBuyer
	}
	err := us.validate.StructCtx(ctx, user)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.Password = string(hashPassword)
	return us.userRepository.Create(ctx, user)
}

func (us *userService) SignIn(ctx context.Context, username, password string) (*domain.User, error) {
	userID, passwordHash, role, err := us.userRepository.GetCredentialsByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return &domain.User{ID: userID, Username: username, Role: role}, nil
}

func (us *userService) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return us.userRepository.GetById(ctx, id)
}

func (us *userService) List(ctx context.Context) ([]domain.User, error) {
	return us.userRepository.List(ctx)
}

func (us *userService) UpdateRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	if !role.IsValid() {
		return domain.ErrInvalidRole
	}
	return us.userRepository.UpdateRole(ctx, id, role)
}

func (us *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return us.userRepository.Delete(ctx, id)
}
