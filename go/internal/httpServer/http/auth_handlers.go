package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
)

type LoginRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secret"`
}

type CurrentUserResponse struct {
	User          models.User `json:"user" binding:"required"`
	Authenticated bool        `json:"authenticated" binding:"required"`
}

// ErrorResponse is a standard error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid username or password"`
}

// @Summary User login
// @Description Authenticates a user and returns their profile. Sets JWT in a secure cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
// @Id Login
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Delegate to the auth service
	result, err := s.Auth.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.Config.AppEnv != "development",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((24 * time.Hour).Seconds()), // adjust TTL
	})

	_ = json.NewEncoder(w).Encode(result.User)
}

// @Summary Dummy User login
// @Description Authenticates a user with the default user in development
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/dummy [post]
// @Id LoginDev
func (s *Server) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if s.Config.AppEnv != "development" {
		http.Error(w, "Not development", http.StatusBadGateway)
	}

	result, err := s.Auth.DummyAuthenticate(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.Config.AppEnv != "development",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((4 * time.Hour).Seconds()),
	})

	_ = json.NewEncoder(w).Encode(result.User)
}

// @Summary CurrentUser
// @Description Gets the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} CurrentUserResponse
// @Router /auth/user [get]
// @Id GetCurrentUser
func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	response := CurrentUserResponse{
		User:          models.User{},
		Authenticated: false,
	}

	c, err := r.Cookie("auth_token")
	if err != nil || c.Value == "" {
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	id, err := s.Auth.Validate(c, s.Config.Secret)
	if err != nil {
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	user, err := s.queries.GetUserByID(r.Context(), id)
	if err != nil {
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response.User = models.MapUser(user)
	response.Authenticated = true

	_ = json.NewEncoder(w).Encode(response)
}

// @Summary Logout current user
// @Description Logs the currently authenticated user out
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/logout [post]
// @Id Logout
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.Config.AppEnv != "development",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,              // Tell browser to delete now
	})
	w.WriteHeader(http.StatusNoContent)
}
