package oidc

import "net/http"

const (
	oidcStateCookie = "oidc_state"
	oidcNonceCookie = "oidc_nonce"
)

type LoginResponse struct {
	Location string        `header:"Location"`
	Cookies  []http.Cookie `header:"Set-Cookie"`
}

type CallbackRequest struct {
	Code        string      `query:"code"`
	State       string      `query:"state"`
	StateCookie http.Cookie `cookie:"oidc_state"`
	NonceCookie http.Cookie `cookie:"oidc_nonce"`
}

type CallbackResponse struct {
	Location string        `header:"Location"`
	Cookies  []http.Cookie `header:"Set-Cookie"`
}
