package service

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	jwtmanager "github.com/K1ver/e-commerce-api/internal/pkg/jwt"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, user *domain.User) (*jwtmanager.TokenPair, error)
	Login(ctx context.Context, username, password string) (*jwtmanager.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*jwtmanager.TokenPair, error)
}

type authService struct {
	userService UserService
	jwt         *jwtmanager.Manager
}

func NewAuthService(userService UserService, jwt *jwtmanager.Manager) AuthService {
	return &authService{userService: userService, jwt: jwt}
}

func (s *authService) Register(ctx context.Context, user *domain.User) (*jwtmanager.TokenPair, error) {
	if err := s.userService.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(user.ID)
}

func (s *authService) Login(ctx context.Context, username, password string) (*jwtmanager.TokenPair, error) {
	userID, err := s.userService.SignIn(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(userID)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*jwtmanager.TokenPair, error) {
	userID, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if _, err := s.userService.GetById(ctx, userID); err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(userID)
}
