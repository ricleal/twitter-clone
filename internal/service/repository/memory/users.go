package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type UserHandler struct {
	Handler
}

func (s *UserHandler) Create(ctx context.Context, u *repository.User) error {
	s.m.Lock()
	defer s.m.Unlock()
	u.ID = uuid.New()
	s.Users = append(s.Users, *u)
	return nil
}

func (s *UserHandler) FindAll(ctx context.Context) ([]repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.Users) == 0 {
		return nil, repository.ErrNotFound
	}
	return s.Users, nil
}

func (s *UserHandler) FindByID(ctx context.Context, id string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.ID == uuid.MustParse(id) {
			return &u, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (s *UserHandler) FindByUsername(ctx context.Context, username string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *UserHandler) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}
