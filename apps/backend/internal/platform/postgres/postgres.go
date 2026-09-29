package postgres

import (
	"database/sql"
	"fmt"

	"nekosync/internal/platform/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Init opens and verifies the PostgreSQL connection pool.
func Init(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
