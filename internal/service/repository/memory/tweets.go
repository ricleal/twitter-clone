package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// TweetHandler is a mock implementation of the repository.TweetServer interface.
type TweetHandler struct {
	Handler
}

// Create creates a new tweet.
func (s *TweetHandler) Create(_ context.Context, t *repository.Tweet) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	t.ID = uuid.New()
	s.Tweets = append(s.Tweets, *t)
	return nil
}

// FindAll returns all tweets.
func (s *TweetHandler) FindAll(_ context.Context) ([]repository.Tweet, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.Tweets, nil
}

// FindByID returns a tweet by ID.
func (s *TweetHandler) FindByID(_ context.Context, id string) (*repository.Tweet, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, t := range s.Tweets {
		if t.ID == uuid.MustParse(id) {
			return &t, nil
		}
	}

	return nil, repository.ErrNotFound
}

// NewTweetHandler returns a new TweetHandler.
func NewTweetHandler() *TweetHandler {
	return &TweetHandler{}
}
