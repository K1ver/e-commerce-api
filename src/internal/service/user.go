package service

import (
	"context"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepository postgres.UserRepository
}

func NewUserService(userRepository postgres.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (us *userService) Create(ctx context.Context, user *domain.User) error {
	return us.userRepository.Create(ctx, user)
}

func (us *userService) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return us.userRepository.GetById(ctx, id)
}

func (us *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return us.userRepository.GetByEmail(ctx, email)
}

func (us *userService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return us.userRepository.GetByUsername(ctx, username)
}

func (us *userService) Update(ctx context.Context, user *domain.User) error {
	return us.userRepository.Update(ctx, user)
}

func (us *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return us.userRepository.Delete(ctx, id)
}
