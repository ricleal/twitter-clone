package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TweetServer struct {
	dbConn repository.DBTx
}

func NewTweetServer(dbConn repository.DBTx) *TweetServer {
	return &TweetServer{
		dbConn: dbConn,
	}
}

func (s *TweetServer) Create(ctx context.Context, t *repository.Tweet) (err error) {

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

func (s *TweetServer) FindAll(ctx context.Context) ([]repository.Tweet, error) {
	ormTweets, err := orm.Tweets().All(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find all tweets: %w", err)
	}

	var tweets []repository.Tweet
	for _, t := range ormTweets {
		tweets = append(tweets, repository.Tweet{
			ID:      uuid.MustParse(t.ID),
			Content: t.Content,
			UserID:  uuid.MustParse(t.UserID),
		})
	}

	return tweets, nil
}

func (s *TweetServer) FindByID(ctx context.Context, id string) (*repository.Tweet, error) {
	ormTweet, err := orm.FindTweet(ctx, s.dbConn, id)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
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
