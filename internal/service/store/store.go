package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
)

type sqlStore struct {
	db DBTx
}

// New creates a new store.
func NewSQLStore(db DBTx) *sqlStore {
	return &sqlStore{
		db: db,
	}
}

func (s *sqlStore) Tweets() repository.TweetRepository {
	return &postgres.TweetServer{}
}

func (s *sqlStore) Users() repository.UserRepository {
	return &postgres.UserServer{}
}

// ExecTx executes the given function within a database transaction.
// we only want to return NewExecTxError for any errors that happen outside the fn function.
func (s *sqlStore) ExecTx(ctx context.Context, fn func(Store) error) error {
	db, ok := s.db.(*sql.DB)
	if !ok {
		return errors.New("ExecTx: db is not a *sql.DB")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	newStore := NewSQLStore(tx)
	err = fn(newStore)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New(err.Error() + "; " + rbErr.Error())
		}
		return err
	}
	return tx.Commit()
}
