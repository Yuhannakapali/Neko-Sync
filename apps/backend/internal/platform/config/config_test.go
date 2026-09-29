package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Run("defaults port and expiry", func(t *testing.T) {
		t.Setenv("PORT", "")
		t.Setenv("JWT_EXPIRY", "")
		t.Setenv("DATABASE_URL", "test_db_url")
		t.Setenv("JWT_SECRET", "secret")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Port != "8080" {
			t.Errorf("expected default port 8080, got %s", cfg.Port)
		}
		if cfg.JWTExpiry != 24*time.Hour {
			t.Errorf("expected default expiry 24h, got %s", cfg.JWTExpiry)
		}
	})

	t.Run("reads values from environment", func(t *testing.T) {
		t.Setenv("PORT", "3000")
		t.Setenv("DATABASE_URL", "test_db_url")
		t.Setenv("JWT_SECRET", "secret")
		t.Setenv("JWT_EXPIRY", "2h")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Port != "3000" || cfg.DatabaseURL != "test_db_url" || cfg.JWTSecret != "secret" || cfg.JWTExpiry != 2*time.Hour {
			t.Errorf("unexpected config: %+v", cfg)
		}
	})

	t.Run("errors without DATABASE_URL", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")
		t.Setenv("JWT_SECRET", "secret")
		if _, err := Load(); err == nil {
			t.Error("expected error for missing DATABASE_URL")
		}
	})

	t.Run("errors without JWT_SECRET", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "test_db_url")
		t.Setenv("JWT_SECRET", "")
		if _, err := Load(); err == nil {
			t.Error("expected error for missing JWT_SECRET")
		}
	})

	t.Run("errors on invalid JWT_EXPIRY", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "test_db_url")
		t.Setenv("JWT_SECRET", "secret")
		t.Setenv("JWT_EXPIRY", "soon")
		if _, err := Load(); err == nil {
			t.Error("expected error for invalid JWT_EXPIRY")
		}
	})
}
