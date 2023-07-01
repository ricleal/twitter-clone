package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func open() (*sql.DB, error) {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	boil.SetDB(db)
	// boil.DebugMode = true
	return db, err
}

type Server struct {
	dbConn *sql.DB
}

func (s *Server) Close() error {
	return s.dbConn.Close()
}

func (s *Server) DB() *sql.DB {
	return s.dbConn
}

func New() *Server {
	db, err := open()
	if err != nil {
		panic(err)
	}
	return &Server{
		dbConn: db,
	}
}
