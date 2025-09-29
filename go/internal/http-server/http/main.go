package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/spa"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	mux *http.ServeMux
	cfg *config.HttpConfig
	sse *sse.SSEManager
}

func NewServer(cfg *config.HttpConfig, pool *pgxpool.Pool) *Server {
	mux := http.NewServeMux()
	huma.DefaultArrayNullable = false
	humaCfg := huma.DefaultConfig("Fogistration", "2.0.0")
	api := humago.New(mux, humaCfg)

	container := container.NewContainer(cfg, pool)

	handlers := handlers.NewHandlers(container)
	handlers.Register(api, "/api")

	container.SSE.CreateEndpoint(api)

	if cfg.AppEnv == "production" {
		distFS, err := spa.SpaFs()
		if err != nil {
			logging.Fatal("failed to get frontend files", err)
		}
		if distFS != nil {
			mux.Handle("/", spaHandler(distFS))
		}
	}

	if err := saveSpec(api); err != nil {
		logging.Error("failed to write api spec", err)
	}

	return &Server{
		mux,
		cfg,
		container.SSE,
	}
}

func (s *Server) Run(ctx context.Context) {
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	go s.sse.Start(ctx, database.GetUrl(&s.cfg.DB))
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

func spaHandler(fs fs.FS) http.Handler {
	fileSrv := http.FileServerFS(fs)
	buildTime := time.Now()

	setCache := func(w http.ResponseWriter, p string) {
		if path.Ext(p) == ".html" {
			w.Header().Set("Cache-Control", "no-cache")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		w.Header().Set("Last-Modified", buildTime.UTC().Format(http.TimeFormat))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := path.Clean(r.URL.Path)

		if path.Ext(clean) != "" {
			setCache(w, clean)
			fileSrv.ServeHTTP(w, r)
			return
		}

		setCache(w, "/")
		r2 := *r
		u := *r.URL
		u.Path = "/"
		r2.URL = &u
		fileSrv.ServeHTTP(w, &r2)
	})
}
