package http

import (
	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/http/middleware"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	netHttp "net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) setupRoutes() {

	s.Router.Use(chiMiddleware.RequestID)
	s.Router.Use(middleware.RequestLogger)
	s.Router.Use(chiMiddleware.Recoverer)

	s.Router.Route("/api", func(r chi.Router) {
		r.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

		r.Get("/health", s.HealthCheck)

		r.Route("/auth", func(auth chi.Router) {
			auth.Post("/login", s.Login)
			auth.Post("/dummy", s.DummyLogin)
			auth.Post("/logout", s.Logout)
			auth.Get("/user", s.GetCurrentUser)
		})

		r.Route("/client", func(client chi.Router) {
			client.Use(s.auth()...)
			client.Get("/", s.ListClients)
			client.Put("/{clientId}/team", s.SetTeam)
		})

		r.Route("/dj", func(dj chi.Router) {
			dj.Use(s.auth()...)
			dj.Get("/team", s.ListTeams)
			dj.Put("/team/{teamId}/client", s.SetTeamClient)
			dj.Get("/contest", s.ListContests)
			dj.Get("/contest/active", s.GetActiveContest)
		})

		r.Route("/wallpaper", func(wallpaper chi.Router) {
			wallpaper.Use(s.auth()...)

			wallpaper.Get("/{contestId}", s.GetWallpaper)
			wallpaper.Put("/{contestId}", s.SetWallpaper)

			wallpaper.Get("/{contestId}/config", s.GetWallpaperConfig)
			wallpaper.Put("/{contestId}/config", s.SetWallpaperConfig)
		})
	})
}

func (s *Server) auth(roles ...string) []func(netHttp.Handler) netHttp.Handler {
	m := []func(netHttp.Handler) netHttp.Handler{
		middleware.Auth(s.Config.Secret, s.Auth),
		middleware.FindUser(s.queries),
	}
	if len(roles) > 0 {
		m = append(m, middleware.RequireRoles(roles...))
	}

	return m
}
