package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type UserServer struct {
	Server
}

func (s *UserServer) Create(ctx context.Context, u repository.User) error {
	u.ID = uuid.New()
	s.Users = append(s.Users, u)
	return nil
}

func (s *UserServer) FindAll(ctx context.Context) ([]repository.User, error) {
	return s.Users, nil
}

func (s *UserServer) FindByID(ctx context.Context, id string) (*repository.User, error) {
	for _, u := range s.Users {
		if u.ID == uuid.MustParse(id) {
			return &u, nil
		}
	}

	return nil, nil
}
