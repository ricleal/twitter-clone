package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

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
	if log.Ctx(ctx).GetLevel() <= zerolog.DebugLevel {
		boil.DebugMode = true
	}
	return db, nil
}

type Handler struct {
	dbConn *sql.DB
}

func (s *Handler) Close() error {
	return s.dbConn.Close()
}

func (s *Handler) DB() *sql.DB {
	return s.dbConn
}

func NewHandler(ctx context.Context) (*Handler, error) {
	db, err := open(ctx)
	if err != nil {
		return nil, err
	}
	log.Ctx(ctx).Info().Int("opened_connections", db.Stats().OpenConnections).
		Msg("Opened database connection")
	return &Handler{
		dbConn: db,
	}, nil
}
