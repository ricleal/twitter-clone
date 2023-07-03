package store

import (
	"context"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// Store is the interface that wraps the repositories.
type Store interface {
	Tweets() repository.TweetRepository
	Users() repository.UserRepository
	ExecTx(ctx context.Context, fn func(Store) error) error
}
