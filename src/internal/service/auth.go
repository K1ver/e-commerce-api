package service

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/domain"
	jwtmanager "github.com/K1ver/e-commerce-api/internal/pkg/jwt"
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
	if user.Role == "" {
		user.Role = domain.RoleBuyer
	}
	if !user.Role.IsValid() {
		return nil, domain.ErrInvalidRole
	}
	if user.Role != domain.RoleBuyer {
		return nil, domain.ErrForbidden
	}
	if err := s.userService.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(jwtmanager.AuthSubject{UserID: user.ID, Role: user.Role})
}

func (s *authService) Login(ctx context.Context, username, password string) (*jwtmanager.TokenPair, error) {
	user, err := s.userService.SignIn(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(jwtmanager.AuthSubject{UserID: user.ID, Role: user.Role})
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*jwtmanager.TokenPair, error) {
	subject, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	user, err := s.userService.GetById(ctx, subject.UserID)
	if err != nil {
		return nil, err
	}
	return s.jwt.GeneratePair(jwtmanager.AuthSubject{UserID: user.ID, Role: user.Role})
}
