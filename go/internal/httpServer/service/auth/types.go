package auth

import (
	"context"
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
)

type AuthResult struct {
	User  models.User
	Token string // JWT for cookie
}

type Service interface {
	Authenticate(ctx context.Context, username, password string) (AuthResult, error)
	DummyAuthenticate(ctx context.Context) (AuthResult, error)
	Validate(cookie *http.Cookie, sec string) (int64, error)
}
