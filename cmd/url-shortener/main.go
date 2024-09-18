package main

import (
	"log/slog"
	"os"

	"github.com/gulmudovv/url-shortener/internal/config"
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
