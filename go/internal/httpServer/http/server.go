package http

import (
	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router  *chi.Mux
	Config  *config.HttpConfig
	queries *database.Queries

	service.ServiceRepo
}

func NewServer(q *database.Queries, cfg *config.HttpConfig) *Server {
	s := &Server{
		Router:      chi.NewRouter(),
		Config:      cfg,
		ServiceRepo: *service.New(q, cfg.Secret),
		queries:     q,
	}

	s.setupRoutes()

	return s
}
