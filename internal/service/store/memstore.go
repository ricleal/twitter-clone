package store

import (
	"context"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

type memStore struct {
	TransactionError bool
	users            repository.UserRepository
	tweets           repository.TweetRepository
}

// NewMemStore creates a new memory store.
func NewMemStore() *memStore {
	// Note that we are using the memory repository implementations here.
	// When we use Tweets() or Users() we are returning the current memory repository
	return &memStore{
		users:  memory.NewUserHandler(),
		tweets: memory.NewTweetHandler(),
	}
}

func (s *memStore) Tweets() repository.TweetRepository {
	return s.tweets
}

func (s *memStore) Users() repository.UserRepository {
	return s.users
}

func (s *memStore) ExecTx(_ context.Context, fn func(Store) error) error {
	if s.TransactionError {
		return NewExecTxError("a transaction related error occurred")
	}
	return fn(s)
}
