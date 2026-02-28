package middleware

import (
	"net/http"
)

// CORS adds the headers that allow the React frontend (running on port 5173)
// to make requests to the Go backend (running on port 8080).
// Without this, the browser blocks cross-origin requests by default.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS is a "preflight" request the browser sends before the real request
		// to check if CORS is allowed. We just return 200 for these.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
