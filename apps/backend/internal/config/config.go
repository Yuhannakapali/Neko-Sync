package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	expiry := 24 * time.Hour
	if raw := os.Getenv("JWT_EXPIRY"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRY %q: %w", raw, err)
		}
		expiry = d
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   secret,
		JWTExpiry:   expiry,
	}, nil
}
