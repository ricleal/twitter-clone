package service

import (
	"context"
	"fmt"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

func validateTweet(content string) error {
	if len(content) > 280 {
		return fmt.Errorf("tweet content is too long")
	}
	return nil
}

type TweetService interface {
	Create(ctx context.Context, t *entities.Tweet) error
	FindAll(ctx context.Context) ([]entities.Tweet, error)
	FindByID(ctx context.Context, id string) (*entities.Tweet, error)
}

type tweetService struct {
	store store.Store
}

func NewTweetService(s store.Store) *tweetService {
	return &tweetService{s}
}

func (s *tweetService) Create(ctx context.Context, t *entities.Tweet) error {
	// open a transaction
	if err := s.store.ExecTx(ctx, func(scopedStore store.Store) error {
		tweetRepo := scopedStore.Tweets()
		userRepo := scopedStore.Users()

		// check if user exists
		_, err := userRepo.FindByID(ctx, t.UserID.String())
		if err != nil {
			if err == repository.ErrNotFound {
				return entities.ErrInvalidUserID
			}
			return fmt.Errorf("error finding user: %w", err)
		}

		// validate tweet
		if err := validateTweet(t.Content); err != nil {
			return fmt.Errorf("invalid tweet: %w", err)
		}

		tweet := &repository.Tweet{
			Content: t.Content,
			UserID:  t.UserID,
		}
		err = tweetRepo.Create(ctx, tweet)
		if err != nil {
			return fmt.Errorf("could not create tweet: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("could not create tweet in the tx: %w", err)
	}
	return nil
}

func (s *tweetService) FindAll(ctx context.Context) ([]entities.Tweet, error) {
	repo := s.store.Tweets()
	tweets, err := repo.FindAll(ctx)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find all tweets: %w", err)
	}

	entTweets := make([]entities.Tweet, 0, len(tweets))
	for _, t := range tweets {
		entTweets = append(entTweets, entities.Tweet{
			ID:      t.ID,
			Content: t.Content,
			UserID:  t.UserID,
		})
	}

	return entTweets, nil
}

func (s *tweetService) FindByID(ctx context.Context, id string) (*entities.Tweet, error) {
	repo := s.store.Tweets()
	t, err := repo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tweet by id %s: %w", id, err)
	}

	if t == nil {
		return nil, nil
	}

	return &entities.Tweet{
		ID:      t.ID,
		Content: t.Content,
		UserID:  t.UserID,
	}, nil
}
