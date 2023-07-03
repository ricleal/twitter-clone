package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// UserHandler is a mock implementation of the repository.UserServer interface.
type UserHandler struct {
	Handler
}

// Create creates a new user.
func (s *UserHandler) Create(_ context.Context, u *repository.User) error {
	s.m.Lock()
	defer s.m.Unlock()
	u.ID = uuid.New()
	s.Users = append(s.Users, *u)
	return nil
}

// FindAll returns all users.
func (s *UserHandler) FindAll(_ context.Context) ([]repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.Users, nil
}

// FindByID returns a user by ID.
func (s *UserHandler) FindByID(_ context.Context, id string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.ID == uuid.MustParse(id) {
			return &u, nil
		}
	}

	return nil, repository.ErrNotFound
}

// FindByUsername returns a user by username.
func (s *UserHandler) FindByUsername(_ context.Context, username string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

// FindByEmail returns a user by email.
func (s *UserHandler) FindByEmail(_ context.Context, email string) (*repository.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.Users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, repository.ErrNotFound
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}
