package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq" // postgres driver registration
	"github.com/stephenafamo/bob"
)

// open opens a connection to the database.
func open(_ context.Context) (*sql.DB, error) {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, errors.New("DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

// Storage is struct that holds the database connection.
type Storage struct {
	dbConn bob.DB
}

// NewStorage returns a new Handler with a database connection.
func NewStorage(ctx context.Context) (*Storage, error) {
	db, err := open(ctx)
	if err != nil {
		return nil, err
	}
	slog.InfoContext( //nolint:sloglint // global logger; slog.SetDefault called before this
		ctx,
		"Opened database connection",
		"opened_connections",
		db.Stats().OpenConnections,
	)
	return &Storage{
		dbConn: bob.NewDB(db),
	}, nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	return s.dbConn.Close() //nolint:wrapcheck //no need to wrap here
}

// DB returns the database connection.
func (s *Storage) DB() bob.DB {
	return s.dbConn
}
