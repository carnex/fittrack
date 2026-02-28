package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/carnex/fittrack/backend/config"
	"github.com/carnex/fittrack/backend/router"
	"github.com/carnex/fittrack/backend/server"
)

func main() {
	// Load .env file in development. In CI/production, env vars are set directly.
	// godotenv.Load() doesn't error if the file is missing — that's intentional.
	_ = godotenv.Load()

	// Load and validate config. If required vars are missing, we fail fast here
	// rather than discovering it later when a request comes in.
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Build the router with all routes registered
	r := router.New()

	// Build the HTTP server
	srv := server.New(cfg.Port, r)
	server.LogStartup(cfg.Port, cfg.Env)

	// Start server in a goroutine so it doesn't block the shutdown logic below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown — wait for Ctrl+C or SIGTERM (what Docker sends on stop)
	// then give in-flight requests 10 seconds to finish before closing.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}

	slog.Info("server stopped")
}
