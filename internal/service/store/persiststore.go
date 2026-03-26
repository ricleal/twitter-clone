package store

import (
	"context"

	"github.com/stephenafamo/bob"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
)

// persistentStore is a store backed by a PostgreSQL database connection.
type persistentStore struct {
	db bob.DB
}

// persistentStoreTx is a transaction-scoped store. ExecTx is a passthrough
// because the underlying driver does not support nested transactions.
type persistentStoreTx struct {
	db bob.Executor
}

// NewPersistentStore creates a new store with the given database connection.
func NewPersistentStore(db bob.DB) *persistentStore {
	return &persistentStore{db: db}
}

// Tweets returns a TweetRepository for managing tweets.
func (s *persistentStore) Tweets() repository.TweetRepository {
	return postgres.NewTweetStorage(s.db)
}

// Users returns a UserRepository for managing users.
func (s *persistentStore) Users() repository.UserRepository {
	return postgres.NewUserStorage(s.db)
}

// ExecTx executes fn within a database transaction.
func (s *persistentStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bob.Executor) error {
		return fn(&persistentStoreTx{db: tx})
	})
}

func (s *persistentStoreTx) Tweets() repository.TweetRepository {
	return postgres.NewTweetStorage(s.db)
}

func (s *persistentStoreTx) Users() repository.UserRepository {
	return postgres.NewUserStorage(s.db)
}

// ExecTx on a transaction-scoped store runs fn directly — nested transactions
// are not supported by the underlying driver.
func (s *persistentStoreTx) ExecTx(_ context.Context, fn func(Store) error) error {
	return fn(s)
}
