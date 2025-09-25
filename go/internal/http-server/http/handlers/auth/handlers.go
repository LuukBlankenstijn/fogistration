package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) login(ctx context.Context, request *loginRequest) (*loginResponse, error) {
	req := request.Body

	result, err := h.S.Auth.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid username or password")
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Cfg.UseHttps,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	}

	return &loginResponse{cookie, result.User}, nil
}

func (h *Handlers) logout(ctx context.Context, request *struct{}) (*logoutResponse, error) {
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Cfg.UseHttps,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}

	return &logoutResponse{cookie}, nil
}

func (h *Handlers) devLogin(ctx context.Context, request *struct{}) (*loginResponse, error) {
	if h.Cfg.AppEnv != "development" {
		return nil, huma.Error404NotFound("")
	}

	result, err := h.S.Auth.DummyAuthenticate(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("")
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Cfg.UseHttps,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	}

	return &loginResponse{cookie, result.User}, nil
}
