package store

import (
	"context"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type Store interface {
	Tweets() repository.TweetRepository
	Users() repository.UserRepository
	ExecTx(ctx context.Context, fn func(Store) error) error
}
