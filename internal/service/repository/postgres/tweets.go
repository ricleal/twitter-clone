package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TweetServer struct {
	Server
}

func (s *TweetServer) Create(ctx context.Context, t service.Tweet) (err error) {
	// start transaction
	tx, err := s.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Make sure the user exists
	_, err = orm.FindUser(ctx, tx, t.UserID.String())
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

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

func (s *TweetServer) FindAll(ctx context.Context) ([]service.Tweet, error) {
	ormTweets, err := orm.Tweets().All(ctx, s.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to find all tweets: %w", err)
	}

	var tweets []service.Tweet
	for _, t := range ormTweets {
		tweets = append(tweets, service.Tweet{
			ID:      uuid.MustParse(t.ID),
			Content: t.Content,
			UserID:  uuid.MustParse(t.UserID),
		})
	}

	return tweets, nil
}

func (s *TweetServer) FindByID(ctx context.Context, id string) (*service.Tweet, error) {
	ormTweet, err := orm.FindTweet(ctx, s.dbConn, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tweet by id: %w", err)
	}

	return &service.Tweet{
		ID:      uuid.MustParse(ormTweet.ID),
		Content: ormTweet.Content,
		UserID:  uuid.MustParse(ormTweet.UserID),
	}, nil
}
