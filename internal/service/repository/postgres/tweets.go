package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/orm"
)

// TweetStorage is a postgres implementation of the repository.TweetStorage interface.
type TweetStorage struct {
	dbConn repository.DBTx
}

// NewTweetStorage returns a new TweetServer.
func NewTweetStorage(dbConn repository.DBTx) *TweetStorage {
	return &TweetStorage{
		dbConn: dbConn,
	}
}

// Create creates a new tweet.
func (s *TweetStorage) Create(ctx context.Context, t *repository.Tweet) (err error) {
	tweet := orm.Tweet{
		ID:      uuid.NewString(),
		Content: t.Content,
		UserID:  t.UserID.String(),
	}

	err = tweet.Insert(ctx, s.dbConn, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert tweet: %w", err)
	}

	return nil
}

// FindAll returns all tweets.
func (s *TweetStorage) FindAll(ctx context.Context) ([]repository.Tweet, error) {
	ormTweets, err := orm.Tweets().All(ctx, s.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to find all tweets: %w", err)
	}

	tweets := make([]repository.Tweet, 0, len(ormTweets))
	for _, t := range ormTweets {
		tweets = append(tweets, repository.Tweet{
			ID:      uuid.MustParse(t.ID),
			Content: t.Content,
			UserID:  uuid.MustParse(t.UserID),
		})
	}

	return tweets, nil
}

// FindByID returns a tweet by ID.
func (s *TweetStorage) FindByID(ctx context.Context, id string) (*repository.Tweet, error) {
	ormTweet, err := orm.FindTweet(ctx, s.dbConn, id)
	if err != nil {
		// Check if the error is a not found error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find tweet by id: %w", err)
	}

	return &repository.Tweet{
		ID:      uuid.MustParse(ormTweet.ID),
		Content: ormTweet.Content,
		UserID:  uuid.MustParse(ormTweet.UserID),
	}, nil
}
