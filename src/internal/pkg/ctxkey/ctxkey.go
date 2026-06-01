package ctxkey

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserID contextKey = "userID"

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, UserID, id)
}

func UserIDFrom(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserID).(uuid.UUID)
	return id, ok
}
