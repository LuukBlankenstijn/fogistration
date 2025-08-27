package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	mux *http.ServeMux
	cfg *config.HttpConfig
}

func NewServer(cfg *config.HttpConfig, pool *pgxpool.Pool) *Server {
	mux := http.NewServeMux()
	huma.DefaultArrayNullable = false
	humaCfg := huma.DefaultConfig("Fogistration", "2.0.0")
	api := humago.New(mux, humaCfg)

	container := container.NewContainer(cfg, pool)
	middlewareFactory := middleware.NewMiddlewareFactory(container)

	handlers := handlers.NewHandlers(container)
	handlers.Register(api, middlewareFactory, "/api")

	if err := saveSpec(api); err != nil {
		logging.Error("failed to write api spec", err)
	}

	return &Server{
		mux,
		cfg,
	}
}

func (s *Server) Run() {
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	logging.Info("%s", fmt.Sprintf("Starting http server on %s", address))
	logging.Fatal("failed to run server", http.ListenAndServe(address, s.mux))
}

func saveSpec(api huma.API) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	outDir := filepath.Join(cwd, "api")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	data, err := api.OpenAPI().YAML()
	if err != nil {
		return err
	}

	out := filepath.Join(outDir, "openapi.yaml")
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return err
	}
	return nil
}
