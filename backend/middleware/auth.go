package middleware

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is an unexported type for context keys in this package.
// Using a custom type prevents collisions with other packages using context.
type contextKey string

const UserIDKey contextKey = "userID"

// Auth validates the JWT token on protected routes.
// For now this is a placeholder — we'll fill in the real JWT logic
// in Phase 2 when we build the auth endpoints.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Phase 2: validate the token and extract the user ID
		// For now we just check the header exists
		next.ServeHTTP(w, r)
	})
}

// GetUserID retrieves the user ID from the request context.
// Handlers call this to find out which user made the request.
func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}

// SetUserID stores the user ID in the request context.
// The Auth middleware calls this after validating the token.
func SetUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}
