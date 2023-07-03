package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

// Create a regular expression to match valid email addresses.
var reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)

func validateEmail(email string) bool {
	// Return true if the email address matches the regular expression.
	return reEmail.MatchString(email)
}

type UserService interface {
	Create(ctx context.Context, u *entities.User) error
	FindAll(ctx context.Context) ([]entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(s store.Store) *userService {
	repo := s.Users()
	return &userService{repo}
}

func (s *userService) Create(ctx context.Context, u *entities.User) error {
	if !validateEmail(u.Email) {
		return entities.ErrInvalidEmail
	}

	user := &repository.User{
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) FindAll(ctx context.Context) ([]entities.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not find users: %w", err)
	}

	entUsers := make([]entities.User, 0, len(users))
	for _, u := range users {
		entUsers = append(entUsers, entities.User{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Name:     u.Name,
		})
	}
	return entUsers, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*entities.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	if u == nil {
		return nil, fmt.Errorf("user with id %s not found", id)
	}

	return &entities.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
	}, nil
}
