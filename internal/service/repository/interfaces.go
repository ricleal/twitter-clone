package repository

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, p *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type TweetRepository interface {
	FindAll(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, t *Tweet) error
	FindByID(ctx context.Context, id string) (*Tweet, error)
}

type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
