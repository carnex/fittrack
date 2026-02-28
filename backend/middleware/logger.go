package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger logs each incoming request with method, path, status, and duration.
// This is useful for debugging and monitoring — you can see every request
// that hits your server in the terminal.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter so we can capture the status code
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.status,
			"duration", time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
// The standard ResponseWriter doesn't expose the status code after it's set,
// so we wrap it to intercept WriteHeader calls.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
