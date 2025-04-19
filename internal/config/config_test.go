package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env.test") // or ".env"
	os.Exit(m.Run())
}
func TestLoad(t *testing.T) {
	t.Run("should load default port if PORT is not set", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Setenv("DATABASE_URL", "test_db_url")
		defer os.Unsetenv("DATABASE_URL")

		config := Load()

		if config.Port != "8080" {
			t.Errorf("expected default port 8080, got %s", config.Port)
		}
	})

	t.Run("should load PORT from environment variable", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("DATABASE_URL", "test_db_url")
		defer os.Unsetenv("PORT")
		defer os.Unsetenv("DATABASE_URL")

		config := Load()

		if config.Port != "3000" {
			t.Errorf("expected port 3000, got %s", config.Port)
		}
	})

	t.Run("should fail if DATABASE_URL is not set", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected log.Fatal to terminate the program, but it did not")
			}
		}()

		Load()
	})

	t.Run("should load DATABASE_URL from environment variable", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "test_db_url")
		defer os.Unsetenv("DATABASE_URL")

		config := Load()

		if config.DatabaseURL != "test_db_url" {
			t.Errorf("expected database URL 'test_db_url', got %s", config.DatabaseURL)
		}
	})
}
