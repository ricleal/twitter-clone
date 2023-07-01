package store

import (
	"context"
	"database/sql"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Store interface {
	Tweets() repository.TweetRepository
	Users() repository.UserRepository
	ExecTx(ctx context.Context, fn func(Store) error) error
}
