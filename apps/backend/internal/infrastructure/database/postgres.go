package database

import (
	"database/sql"
	"log"

	"nekosync/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Init initializes the PostgreSQL database connection
func Init(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db
}
