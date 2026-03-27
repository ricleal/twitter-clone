package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/models"
)

// TweetStorage is a postgres implementation of the repository.TweetStorage interface.
type TweetStorage struct {
	dbConn bob.Executor
}

// NewTweetStorage returns a new TweetServer.
func NewTweetStorage(dbConn bob.Executor) *TweetStorage {
	return &TweetStorage{
		dbConn: dbConn,
	}
}

// Create creates a new tweet.
func (s *TweetStorage) Create(ctx context.Context, t *entities.Tweet) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate tweet id: %w", err)
	}
	setter := &models.TweetSetter{
		ID:      omit.From(id.String()),
		Content: omit.From(t.Content),
		UserID:  omit.From(t.UserID.String()),
	}

	_, err = models.Tweets.Insert(setter).One(ctx, s.dbConn)
	if err != nil {
		return fmt.Errorf("failed to insert tweet: %w", err)
	}

	return nil
}

// FindAll returns all tweets.
func (s *TweetStorage) FindAll(ctx context.Context) ([]entities.Tweet, error) {
	ormTweets, err := models.Tweets.Query().All(ctx, s.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to find all tweets: %w", err)
	}

	tweets := make([]entities.Tweet, 0, len(ormTweets))
	for _, t := range ormTweets {
		tweets = append(tweets, entities.Tweet{
			ID:      uuid.MustParse(t.ID),
			Content: t.Content,
			UserID:  uuid.MustParse(t.UserID),
		})
	}

	return tweets, nil
}

// FindByID returns a tweet by ID.
func (s *TweetStorage) FindByID(ctx context.Context, id string) (*entities.Tweet, error) {
	ormTweet, err := models.Tweets.Query(
		sm.Where(models.Tweets.Columns.ID.EQ(psql.Arg(id))),
	).One(ctx, s.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find tweet by id: %w", err)
	}

	return &entities.Tweet{
		ID:      uuid.MustParse(ormTweet.ID),
		Content: ormTweet.Content,
		UserID:  uuid.MustParse(ormTweet.UserID),
	}, nil
}
