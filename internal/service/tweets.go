package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

// validateTweet validates a tweet content is not too long.
func validateTweet(content string) error {
	if len(content) > 280 {
		return fmt.Errorf("tweet content is too long")
	}
	return nil
}

// TweetService is a domain service for tweets.
type TweetService interface {
	Create(ctx context.Context, t *entities.Tweet) error
	FindAll(ctx context.Context) ([]entities.Tweet, error)
	FindByID(ctx context.Context, id string) (*entities.Tweet, error)
}

// tweetService is an implementation of the TweetService interface.
type tweetService struct {
	store store.Store
}

// NewTweetService returns a new TweetService.
func NewTweetService(s store.Store) *tweetService {
	return &tweetService{s}
}

// Create creates a new tweet.
func (s *tweetService) Create(ctx context.Context, t *entities.Tweet) error {
	// open a transaction
	if errOut := s.store.ExecTx(ctx, func(scopedStore store.Store) error {
		tweetRepo := scopedStore.Tweets()
		userRepo := scopedStore.Users()

		// check if user exists
		_, err := userRepo.FindByID(ctx, t.UserID.String())
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return entities.ErrInvalidUserID
			}
			return fmt.Errorf("error finding user: %w", err)
		}

		// validate tweet
		if err = validateTweet(t.Content); err != nil {
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
	}); errOut != nil {
		return fmt.Errorf("could not create tweet in the tx: %w", errOut)
	}
	return nil
}

// FindAll returns all tweets.
func (s *tweetService) FindAll(ctx context.Context) ([]entities.Tweet, error) {
	repo := s.store.Tweets()
	tweets, err := repo.FindAll(ctx)
	if err != nil {
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

// FindByID returns a tweet by ID.
func (s *tweetService) FindByID(ctx context.Context, id string) (*entities.Tweet, error) {
	repo := s.store.Tweets()
	t, err := repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
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
