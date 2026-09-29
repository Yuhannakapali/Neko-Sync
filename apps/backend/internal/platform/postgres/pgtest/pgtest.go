// Package pgtest opens the integration-test database. Tests using it are
// skipped unless TEST_DATABASE_URL points at a migrated Postgres; every call
// empties the application tables, so never point it at real data.
package pgtest

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(`TRUNCATE users, works CASCADE`); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	return db
}
