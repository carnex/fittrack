package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/carnex/fittrack/backend/config"
	"github.com/carnex/fittrack/backend/store"
)

type AppData struct {
	Config *config.Config
	Store  store.Store
}

func New(app *AppData, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", app.Config.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func LogStartup(port, env string) {
	slog.Info("server starting",
		"port", port,
		"env", env,
		"url", fmt.Sprintf("http://localhost:%s", port),
	)
}
