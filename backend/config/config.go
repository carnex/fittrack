package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
// Using a struct means config is explicit and type-safe —
// no scattered os.Getenv() calls throughout the codebase.
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	Env         string
}

// Load reads environment variables and returns a Config.
// We call this once at startup in main.go and pass it around.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // sensible default
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
		JWTSecret:   jwtSecret,
		Env:         os.Getenv("ENV"),
	}, nil
}
