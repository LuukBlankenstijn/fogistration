package main

import (
	"fmt"
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/docs"
	httpServer "github.com/LuukBlankenstijn/fogistration/internal/http"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Fogistration server API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	var cfg config.HttpConfig
	err := config.Load(&cfg, ".env-http")
	if err != nil {
		logging.Fatal("Failed to load config: %v", err)
	}

	logging.SetupLogger(cfg.LogLevel, cfg.AppEnv)

	// test
	// url := database.GetUrl(&cfg)
	// dbpool, err := pgxpool.New(ctx, url)
	// if err != nil {
	// 	logging.Fatal("unable to create dbpool: %w", err)
	// }
	// defer dbpool.Close()
	//
	// queries := database.New(dbpool)

	httpServer := httpServer.NewServer(&cfg)

	initHttpSwagger(httpServer.Router)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	logging.Info("Starting server on %s", addr)
	logging.Info("Swagger available at http://%s/swagger/index.html", addr)
	if err := http.ListenAndServe(addr, httpServer.Router); err != nil {
		logging.Fatal("Server failed: %v", err)
	}

}

func initHttpSwagger(router *chi.Mux) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// This is a workaround to serve the swagger doc.json file, there was an issue with swagger not serving it correctly
	router.Get("/swagger/doc.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})
}
