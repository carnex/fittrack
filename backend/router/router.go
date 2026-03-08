package router

import (
	"encoding/json"
	"net/http"

	handlers "github.com/carnex/fittrack/backend/handlers"
	"github.com/carnex/fittrack/backend/server"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(app *server.AppData) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	userHandler := handlers.NewUserHandler(app.UserService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.Register)

		})
	})

	return r
}
