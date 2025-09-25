package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) handleLogin(ctx context.Context, _ *struct{}) (*LoginResponse, error) {
	state := randStr(24)
	nonce := randStr(24)

	u, err := h.S.OIDCService.AuthURL(ctx, state, nonce)
	if err != nil {
		logging.Error("OIDC build auth URL failed", err)
		return nil, huma.Error500InternalServerError("oidc auth url")
	}

	return &LoginResponse{
		Location: u,
		Cookies: []http.Cookie{
			http.Cookie{
				Name:     oidcStateCookie,
				Value:    state,
				Path:     "/",
				MaxAge:   300,
				HttpOnly: true,
				Secure:   h.Cfg.UseHttps,
				SameSite: http.SameSiteLaxMode,
			},
			http.Cookie{
				Name:     oidcNonceCookie,
				Value:    nonce,
				Path:     "/",
				MaxAge:   300,
				HttpOnly: true,
				Secure:   h.Cfg.UseHttps,
				SameSite: http.SameSiteLaxMode,
			}},
	}, nil
}

func (h *Handlers) handleCallback(ctx context.Context, req *CallbackRequest) (*CallbackResponse, error) {
	res, err := h.S.OIDCService.HandleCallback(ctx, req.Code, req.NonceCookie.Value)
	if err != nil {
		logging.Error("OIDC callback failed", err)
		return nil, huma.Error401Unauthorized("oidc verify failed")
	}

	redirectTo := "/"

	return &CallbackResponse{
		Location: redirectTo,
		Cookies: []http.Cookie{
			http.Cookie{
				Name:     "auth_token",
				Value:    res.Token,
				Path:     "/",
				HttpOnly: true,
				Secure:   h.Cfg.UseHttps,
				SameSite: http.SameSiteLaxMode,
			},
			http.Cookie{
				Name:    oidcStateCookie,
				Value:   "",
				Path:    "/",
				MaxAge:  -1,
				Expires: time.Unix(0, 0),
			},
			http.Cookie{
				Name:    oidcNonceCookie,
				Value:   "",
				Path:    "/",
				MaxAge:  -1,
				Expires: time.Unix(0, 0),
			},
		},
	}, nil
}

func (h *Handlers) handleIsEnabled(ctx context.Context, _ *struct{}) (*struct{}, error) {
	return &struct{}{}, nil
}
func randStr(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
