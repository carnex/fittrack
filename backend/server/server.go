package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// New creates a configured http.Server.
// We set timeouts explicitly — without them, slow clients
// can hold connections open forever and exhaust resources.
func New(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// LogStartup prints a friendly startup message.
func LogStartup(port, env string) {
	slog.Info("server starting",
		"port", port,
		"env", env,
		"url", fmt.Sprintf("http://localhost:%s", port),
	)
}
