package store

import (
	"context"

	memdb "github.com/hashicorp/go-memdb"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

type memStore struct {
	db               *memdb.MemDB
	TransactionError bool
}

// NewMemStore creates a new memory store backed by go-memdb.
func NewMemStore() *memStore { //nolint:revive // returns unexported type; tests access TransactionError field directly
	db, err := memory.NewDB()
	if err != nil {
		panic("failed to create in-memory database: " + err.Error())
	}
	return &memStore{db: db}
}

func (s *memStore) Tweets() repository.TweetRepository {
	return memory.NewTweetHandler(s.db)
}

func (s *memStore) Users() repository.UserRepository {
	return memory.NewUserHandler(s.db)
}

func (s *memStore) ExecTx(_ context.Context, fn func(Store) error) error {
	if s.TransactionError {
		return NewExecTxError("a transaction related error occurred")
	}
	return fn(s)
}
