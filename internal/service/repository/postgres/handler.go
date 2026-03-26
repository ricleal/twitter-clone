package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// open opens a connection to the database.
func open(ctx context.Context) (*sql.DB, error) {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	boil.SetDB(db)
	if slog.Default().Enabled(ctx, slog.LevelDebug-4) {
		boil.DebugMode = true
	}
	return db, nil
}

// Storage is struct that holds the database connection.
type Storage struct {
	dbConn *sql.DB
}

// Close closes the database connection.
func (s *Storage) Close() error {
	return s.dbConn.Close() //nolint:wrapcheck //no need to wrap here
}

// DB returns the database connection.
func (s *Storage) DB() *sql.DB {
	return s.dbConn
}

// NewStorage returns a new Handler with a database connection.
func NewStorage(ctx context.Context) (*Storage, error) {
	db, err := open(ctx)
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "Opened database connection", "opened_connections", db.Stats().OpenConnections)
	return &Storage{
		dbConn: db,
	}, nil
}
