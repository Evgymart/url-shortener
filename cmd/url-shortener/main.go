package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"urlshort/internal/config"
	del "urlshort/internal/http-server/handlers/url/delete"
	"urlshort/internal/http-server/handlers/url/redirect"
	"urlshort/internal/http-server/handlers/url/save"
	"urlshort/internal/storage/dragonfly"
	middle "urlshort/lib/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middle.LoggerMiddleware(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.Auth.User: cfg.Auth.Password,
		}))

		r.Post("/url", save.New(logger, storage))
		r.Delete("/url", del.New(logger, storage))
	})

	router.Get("/{alias}", redirect.New(logger, storage))
	logger.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error("Server error: " + err.Error())
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
