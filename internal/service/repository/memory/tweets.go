package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type TweetHandler struct {
	Handler
}

func (s *TweetHandler) Create(ctx context.Context, t *repository.Tweet) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	t.ID = uuid.New()
	s.Tweets = append(s.Tweets, *t)
	return nil
}

func (s *TweetHandler) FindAll(ctx context.Context) ([]repository.Tweet, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.Tweets) == 0 {
		return nil, repository.ErrNotFound
	}
	return s.Tweets, nil
}

func (s *TweetHandler) FindByID(ctx context.Context, id string) (*repository.Tweet, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, t := range s.Tweets {
		if t.ID == uuid.MustParse(id) {
			return &t, nil
		}
	}

	return nil, repository.ErrNotFound
}

func NewTweetHandler() *TweetHandler {
	return &TweetHandler{}
}
