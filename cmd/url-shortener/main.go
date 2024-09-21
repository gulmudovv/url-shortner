package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gulmudovv/url-shortener/internal/config"
	"github.com/gulmudovv/url-shortener/internal/http-server/handlers/redirect"
	"github.com/gulmudovv/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/gulmudovv/url-shortener/internal/http-server/middleware/logger"
	"github.com/gulmudovv/url-shortener/internal/lib/logger/sl"
	"github.com/gulmudovv/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//init config
	cfg := config.MustLoad()

	//init logger
	logger := setupLogger(cfg.Env)
	logger.Info("url shortener", slog.String("env", cfg.Env))
	logger.Debug("debug message working")

	//init database

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// init router
	router := chi.NewRouter()

	//middelware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/url", func(r chi.Router) {

		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(logger, storage))

		//TODO delete handler
	})

	//routes

	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	//server

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
