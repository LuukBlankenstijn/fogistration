package http

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router *chi.Mux
	Config *config.HttpConfig
}

func NewServer(cfg *config.HttpConfig) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		Config: cfg,
	}

	s.setupRoutes()

	return s
}
