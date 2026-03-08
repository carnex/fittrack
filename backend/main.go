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
	"github.com/carnex/fittrack/backend/service"
	"github.com/carnex/fittrack/backend/store"
)

func main() {

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dbConn, err := store.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to establish DB Connection", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close()
	pgStore := store.NewPostgresStore(dbConn)
	userService := service.NewUserService(pgStore)
	appdata := server.AppData{
		Config:      cfg,
		Store:       pgStore,
		UserService: userService,
	}

	r := router.New(&appdata)

	srv := server.New(&appdata, r)
	server.LogStartup(cfg.Port, cfg.Env)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}

	slog.Info("server stopped")
}
