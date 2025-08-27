package auth

import (
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
)

type loginRequest struct {
	Body struct {
		Username string `json:"username" doc:"Username"`
		Password string `json:"password" doc:"Password"`
	}
}

type loginResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      models.User
}

type logoutResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}
