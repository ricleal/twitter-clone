package repository

import (
	"context"
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
