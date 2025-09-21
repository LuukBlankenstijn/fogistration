package middleware

import (
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
)

const AUTH_COOKIE = "auth_token"

// Auth validates the "auth_token" cookie and injects the user ID into context.
func (m *MiddlewareFactory) ValidateAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		c, err := huma.ReadCookie(ctx, AUTH_COOKIE)
		if c == nil || c.Value == "" || err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing auth cookie")
			return
		}

		id, err := m.S.Auth.Validate(c, m.Cfg.Secret)
		if err != nil {
			logging.Error("Auth middleware: failed to validate token", err)
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx = huma.WithValue(ctx, userIDCtxKey, id)
		next(ctx)
	}
}

// FindUser loads the full user from DB using the ID placed by Auth.
func (m *MiddlewareFactory) FindUser(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		id, ok := GetUserID(ctx.Context())
		if !ok {
			logging.Error("FindUser middleware: user ID not found in context", nil)
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "user id not set")
			return
		}

		user, err := m.Q.GetUserByID(ctx.Context(), id)
		if err != nil {
			logging.Error("FindUser middleware: user lookup failed", err, "userID", id)
			_ = huma.WriteErr(api, ctx, http.StatusBadRequest, "user not found")
			return
		}

		ctx = huma.WithValue(ctx, userCtxKey, user)
		next(ctx)
	}
}

// RequireRoles ensures the authenticated user has one of the allowed roles.
func (m *MiddlewareFactory) RequireRoles(api huma.API, roles ...models.UserRole) func(huma.Context, func(huma.Context)) {
	allowed := map[models.UserRole]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(ctx huma.Context, next func(huma.Context)) {
		u, ok := ctx.Context().Value(userCtxKey).(database.User)
		if !ok {
			logging.Error("RequireRoles middleware: user not loaded in context", nil)
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "user not loaded")
			return
		}

		if len(allowed) > 0 {
			if _, ok := allowed[u.Role]; !ok {
				logging.Warn("RequireRoles middleware: forbidden role: userId: %d, role: %s, required: %+v", u.ID, u.Role, roles)
				_ = huma.WriteErr(api, ctx, http.StatusForbidden, "forbidden")
				return
			}
		}

		next(ctx)
	}
}

// Stacks the VerifyAuth, FindUser, and VerifyAuth middleware
func (m *MiddlewareFactory) Auth(api huma.API, roles ...models.UserRole) huma.Middlewares {
	middleware := huma.Middlewares{
		m.ValidateAuth(api),
		m.FindUser(api),
	}
	if len(roles) > 0 {
		middleware = append(middleware, m.RequireRoles(api, roles...))
	}
	return middleware
}
