package db

import (
	"database/sql"
	"log"

	"nekosync/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Init(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
