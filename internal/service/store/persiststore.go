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
type sqlStore struct {
	db repository.DBTx
}

// NewSQLStore creates a new store with the given database connection.
func NewSQLStore(db repository.DBTx) *sqlStore {
	return &sqlStore{
		db: db,
	}
}

// Tweets returns a TweetRepository for managing tweets.
func (s *sqlStore) Tweets() repository.TweetRepository {
	return postgres.NewTweetStorage(s.db)
}

// Users returns a UserRepository for managing users.
func (s *sqlStore) Users() repository.UserRepository {
	return postgres.NewUserStorage(s.db)
}

// ExecTx executes the given function within a database transaction.
func (s *sqlStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	db, ok := s.db.(*sql.DB)
	if !ok {
		return errors.New("ExecTx: db is not a *sql.DB")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	newStore := NewSQLStore(tx)
	err = fn(newStore)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("ExecTx: %w: Rollback: %w", err, rbErr)
		}
		return fmt.Errorf("ExecTx: %w", err)
	}
	return tx.Commit() //nolint:wrapcheck //no need to wrap here
}
