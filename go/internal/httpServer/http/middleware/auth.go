package middleware

import (
	"context"
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/service/auth"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

type ctxKey string

const userIdCtxKey ctxKey = "auth.user.id"
const userCtxKey ctxKey = "auth.user"

// Auth validates JWT from the "auth_token" cookie and injects user id into context.
func Auth(secret string, authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("auth_token")
			if err != nil || c.Value == "" {
				http.Error(w, "missing auth cookie", http.StatusUnauthorized)
				return
			}

			id, err := authService.Validate(c, secret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIdCtxKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles ensures the authenticated user has one of the roles.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := GetUser(r.Context())
			if len(allowed) > 0 {
				if _, ok := allowed[u.Role]; !ok {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func FindUser(q *database.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := GetUserId(r.Context())
			if !ok {
				http.Error(w, "user id not set", http.StatusInternalServerError)
				return
			}
			user, err := q.GetUserByID(r.Context(), id)
			if err != nil {
				http.Error(w, "User not found", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}

// GetUserId pulls the identity from context.
func GetUserId(ctx context.Context) (int64, bool) {
	u, ok := ctx.Value(userIdCtxKey).(int64)
	return u, ok
}

func GetUser(ctx context.Context) database.User {
	u, ok := ctx.Value(userCtxKey).(database.User)
	if !ok {
		logging.Fatal("Used GetUser on a context without user", nil)
	}
	return u
}
