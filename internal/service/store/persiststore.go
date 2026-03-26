package store

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
)

// Store is a store for tweets and users.
type persistentStore struct {
	db bob.Executor
}

// NewPersistentStore creates a new store with the given database connection.
func NewPersistentStore(db bob.Executor) *persistentStore {
	return &persistentStore{
		db: db,
	}
}

// Tweets returns a TweetRepository for managing tweets.
func (s *persistentStore) Tweets() repository.TweetRepository {
	return postgres.NewTweetStorage(s.db)
}

// Users returns a UserRepository for managing users.
func (s *persistentStore) Users() repository.UserRepository {
	return postgres.NewUserStorage(s.db)
}

// ExecTx executes the given function within a database transaction.
func (s *persistentStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	db, ok := s.db.(bob.DB)
	if !ok {
		return fmt.Errorf("ExecTx: db is not a bob.DB")
	}
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bob.Executor) error {
		return fn(NewPersistentStore(tx))
	})
}
