package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// SetupTestDB creates a test database connection for testing
func SetupTestDB(t *testing.T) *sql.DB {
	// Use test database connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5433/tall_affiliate?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}
