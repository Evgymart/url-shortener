package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
	"urlshort/internal/config"
	"urlshort/internal/storage/dragonfly"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)
	if logger == nil {
		panic("logger not set")
	}

	storage := dragonfly.New(cfg.Dragonfly)
	if storage == nil {
		panic("storage not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := storage.Ping(ctx); err == nil {
		logger.Info("Pong")
	} else {
		panic("storage error " + err.Error())
	}

	fmt.Println(cfg)
	logger.Info("starting server", slog.String("env", cfg.Env))
	logger.Debug("debug logs are enabled")
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
