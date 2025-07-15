package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/sl"
	"urlShortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting app")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init db", sl.Err(err))
		return
	}

	err = storage.DeleteUrl("google")
	if err != nil {
		log.Error("failed to del url", sl.Err(err))
		return
	}

	//TODO: init router chi
	//TODO: run server

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
