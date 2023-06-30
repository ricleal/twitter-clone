package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func open() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	boil.SetDB(db)
	return db, err
}

type Server struct {
	dbConn *sql.DB
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
