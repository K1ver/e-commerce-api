package jwt

import (
	"fmt"
	"time"

	"github.com/K1ver/e-commerce-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type Claims struct {
	UserID    uuid.UUID `json:"userId"`
	TokenType string    `json:"tokenType"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		accessSecret:  []byte(cfg.Secret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTokenExpireDuration,
		refreshTTL:    cfg.RefreshTokenExpireDuration,
	}
}

func (m *Manager) GeneratePair(userID uuid.UUID) (*TokenPair, error) {
	access, err := m.sign(userID, tokenTypeAccess, m.accessSecret, m.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}
	refresh, err := m.sign(userID, tokenTypeRefresh, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

func (m *Manager) ParseAccess(token string) (uuid.UUID, error) {
	return m.parse(token, tokenTypeAccess, m.accessSecret)
}

func (m *Manager) ParseRefresh(token string) (uuid.UUID, error) {
	return m.parse(token, tokenTypeRefresh, m.refreshSecret)
}

func (m *Manager) sign(userID uuid.UUID, tokenType string, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func (m *Manager) parse(tokenString, expectedType string, secret []byte) (uuid.UUID, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w", err)
	}
	if !t.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if claims.TokenType != expectedType {
		return uuid.Nil, fmt.Errorf("invalid token type")
	}
	return claims.UserID, nil
}
