package middleware

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

type ctxKey string

const (
	userIDCtxKey ctxKey = "auth.user.id"
	userCtxKey   ctxKey = "auth.user"
)

func GetUserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDCtxKey).(int64)
	if !ok {
		logging.Warn("GetUserID: no user ID in context")
		return 0, false
	}
	return v, true
}

func User(ctx context.Context) (database.User, bool) {
	u, ok := ctx.Value(userCtxKey).(database.User)
	if !ok {
		logging.Warn("User: no user object in context")
		return database.User{}, false
	}
	return u, true
}
