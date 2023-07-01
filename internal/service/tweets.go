package service

import (
	"context"
	"fmt"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

type TweetService interface {
	Create(ctx context.Context, t entities.Tweet) error
	FindAll(ctx context.Context) ([]entities.Tweet, error)
	FindByID(ctx context.Context, id string) (*entities.Tweet, error)
}

type tweetService struct {
	store store.Store
}

func NewTweetService(store store.Store) *tweetService {
	return &tweetService{store}
}

func (s *tweetService) Create(ctx context.Context, t entities.Tweet) error {

	// open a transaction
	if err := s.store.ExecTx(ctx, func(scopedStore store.Store) error {
		tweetRepo := scopedStore.Tweets()
		userRepo := scopedStore.Users()

		// check if user exists
		user, err := userRepo.FindByID(ctx, t.UserID.String())
		if err != nil || user == nil {
			return fmt.Errorf("user with id %s not found: %w", t.UserID.String(), err)
		}

		tweet := repository.Tweet{
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
		return nil, err
	}

	var entTweets []entities.Tweet
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
		return nil, err
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
