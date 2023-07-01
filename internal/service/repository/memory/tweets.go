package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type TweetServer struct {
	Server
}

func (s *TweetServer) Create(ctx context.Context, t repository.Tweet) (err error) {
	t.ID = uuid.New()
	s.Tweets = append(s.Tweets, t)
	return nil
}

func (s *TweetServer) FindAll(ctx context.Context) ([]repository.Tweet, error) {
	return s.Tweets, nil
}

func (s *TweetServer) FindByID(ctx context.Context, id string) (*repository.Tweet, error) {
	for _, t := range s.Tweets {
		if t.ID == uuid.MustParse(id) {
			return &t, nil
		}
	}

	return nil, nil
}
