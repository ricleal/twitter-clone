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
	db     *sql.DB
	dbConn bob.DB
	logger *slog.Logger
}

// NewStorage returns a new Handler with a database connection.
func NewStorage(ctx context.Context, logger *slog.Logger) (*Storage, error) {
	log := logger.With("component", "postgres")
	db, err := open(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("Opened database connection", "opened_connections", db.Stats().OpenConnections)
	if pingErr := db.PingContext(ctx); pingErr != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to verify database connection: %w", pingErr)
	}
	return &Storage{
		db:     db,
		dbConn: bob.NewDB(db),
		logger: log,
	}, nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	return s.dbConn.Close() //nolint:wrapcheck //no need to wrap here
}

// Ping verifies the database connection is still alive.
func (s *Storage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx) //nolint:wrapcheck // caller decides how to handle the error
}

// DB returns the database connection.
func (s *Storage) DB() bob.DB {
	return s.dbConn
}
