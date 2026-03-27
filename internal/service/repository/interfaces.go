package repository

import (
	"context"

	"github.com/ricleal/twitter-clone/internal/entities"
)

// UserRepository represents a repository for users.
type UserRepository interface {
	FindAll(ctx context.Context) ([]entities.User, error)
	Create(ctx context.Context, u *entities.User) error
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}

// TweetRepository represents a repository for tweets.
type TweetRepository interface {
	FindAll(ctx context.Context) ([]entities.Tweet, error)
	Create(ctx context.Context, t *entities.Tweet) error
	FindByID(ctx context.Context, id string) (*entities.Tweet, error)
}
