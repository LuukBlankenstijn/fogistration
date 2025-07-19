package http

import (
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func (s *Server) setupRoutes() {

	s.Router.Use(middleware.RequestID)
	s.Router.Use(RequestLogger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))

		r.Get("/health", s.HealthCheck)
	})
}
