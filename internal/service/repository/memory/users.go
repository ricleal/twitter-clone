package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service"
)

type UserServer struct {
	Server
}

func (s *UserServer) Create(ctx context.Context, u service.User) error {
	u.ID = uuid.New()
	s.Users = append(s.Users, u)
	return nil
}

func (s *UserServer) FindAll(ctx context.Context) ([]service.User, error) {
	return s.Users, nil
}

func (s *UserServer) FindByID(ctx context.Context, id string) (*service.User, error) {
	for _, u := range s.Users {
		if u.ID == uuid.MustParse(id) {
			return &u, nil
		}
	}

	return nil, nil
}
