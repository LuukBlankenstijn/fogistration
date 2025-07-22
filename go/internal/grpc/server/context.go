package server

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type contextKey string

const clientContextKey contextKey = "client"

func withClient(ctx context.Context, client database.Client) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}

func getClient(ctx context.Context) (database.Client, bool) {
	client, ok := ctx.Value(clientContextKey).(database.Client)
	return client, ok
}
