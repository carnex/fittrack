package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// New builds and returns the application router.
// All routes are registered here — one place to see the full API surface.
func New() http.Handler {
	r := chi.NewRouter()

	// Global middleware — runs on every request
	r.Use(chimiddleware.RequestID) // Adds a unique ID to each request (useful for logs)
	r.Use(chimiddleware.Recoverer) // Catches panics so the server doesn't crash

	// Health check — used by Docker and CI to confirm the server is alive
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API routes — we'll add these as we build each domain
	r.Route("/api", func(r chi.Router) {
		// auth, workouts, metrics etc. will be registered here
	})

	return r
}
