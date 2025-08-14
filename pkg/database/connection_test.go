package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	t.Run("creates connection with config", func(t *testing.T) {
		config := &Config{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "postgres",
			Database:        "testdb",
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		}

		conn, err := NewConnection(config)
		if err != nil {
			t.Skip("Database not available")
		}
		defer conn.Close()

		// Test connection
		ctx := context.Background()
		err = conn.PingContext(ctx)
		assert.NoError(t, err)

		// Check pool settings
		stats := conn.Stats()
		assert.LessOrEqual(t, stats.MaxOpenConnections, 10)
	})

	t.Run("validates config", func(t *testing.T) {
		config := &Config{
			Host: "",
			Port: 0,
		}

		_, err := NewConnection(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})
}

func TestExecuteInTransaction(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("commits on success", func(t *testing.T) {
		var result int
		err := ExecuteInTransaction(ctx, db, func(tx *sql.Tx) error {
			// Simple query that should succeed
			return tx.QueryRow("SELECT 1").Scan(&result)
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("rollback on error", func(t *testing.T) {
		err := ExecuteInTransaction(ctx, db, func(tx *sql.Tx) error {
			// This should fail
			_, err := tx.Exec("INSERT INTO non_existent_table (id) VALUES (1)")
			return err
		})

		assert.Error(t, err)
	})

	t.Run("handles panic", func(t *testing.T) {
		err := ExecuteInTransaction(ctx, db, func(tx *sql.Tx) error {
			panic("test panic")
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic")
	})
}

func TestWithRetry(t *testing.T) {
	t.Run("succeeds on first try", func(t *testing.T) {
		attempts := 0
		err := WithRetry(3, time.Millisecond, func() error {
			attempts++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("retries on failure", func(t *testing.T) {
		attempts := 0
		err := WithRetry(3, time.Millisecond, func() error {
			attempts++
			if attempts < 3 {
				return assert.AnError
			}
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("fails after max retries", func(t *testing.T) {
		attempts := 0
		err := WithRetry(3, time.Millisecond, func() error {
			attempts++
			return assert.AnError
		})

		assert.Error(t, err)
		assert.Equal(t, 3, attempts)
	})
}

func getTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Skip("Database not available")
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Skip("Database not available")
	}

	return db
}