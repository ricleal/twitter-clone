package repository

import (
	"context"
	"database/sql"
)

// UserRepository represents a repository for users.
type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, p *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// TweetRepository represents a repository for tweets.
type TweetRepository interface {
	FindAll(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, t *Tweet) error
	FindByID(ctx context.Context, id string) (*Tweet, error)
}

// DBTx represents a database transaction or connection interface.
type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
