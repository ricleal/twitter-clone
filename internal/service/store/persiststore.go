package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
)

// Store is a store for tweets and users.
type persistentStore struct {
	db repository.DBTx
}

// NewPersistentStore creates a new store with the given database connection.
func NewPersistentStore(db repository.DBTx) *persistentStore {
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
	db, ok := s.db.(*sql.DB)
	if !ok {
		return errors.New("ExecTx: db is not a *sql.DB")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	newStore := NewPersistentStore(tx)
	err = fn(newStore)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("ExecTx: %s: Rollback: %w", err.Error(), rbErr)
		}
		return fmt.Errorf("ExecTx: %w", err)
	}
	return tx.Commit() //nolint:wrapcheck //no need to wrap here
}
