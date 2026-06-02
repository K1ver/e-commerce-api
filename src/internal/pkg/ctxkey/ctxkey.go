package ctxkey

import (
	"context"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

func WithAuth(ctx context.Context, userID uuid.UUID, role domain.Role) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	return context.WithValue(ctx, roleKey, role)
}

func UserIDFrom(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func RoleFrom(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(roleKey).(domain.Role)
	return role, ok
}
