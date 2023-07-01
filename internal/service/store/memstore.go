package store

import (
	"context"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
)

type memStore struct {
	TransactionError bool
}

// New creates a new store.
func NewMemStore() *memStore {
	return &memStore{}
}

func (s *memStore) Tweets() repository.TweetRepository {
	return &postgres.TweetServer{}
}

func (s *memStore) Users() repository.UserRepository {
	return &postgres.UserServer{}
}

func (s *memStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	if s.TransactionError {
		return NewExecTxError("a transaction related error occurred")
	}
	return fn(s)
}
