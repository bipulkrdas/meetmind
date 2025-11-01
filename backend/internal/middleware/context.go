package middleware

import (
	"context"
	"livekit-consulting/backend/internal/model"
)

// A private key type to prevent collisions
type contextKey string

const userContextKey = contextKey("user")

// WithUser adds the user to the context
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFrom extracts the user from the context
func UserFrom(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}
