package memory

import (
	"context"
	"fmt"

	memdb "github.com/hashicorp/go-memdb"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// UserHandler is a memory implementation of the repository.UserRepository interface.
type UserHandler struct {
	db *memdb.MemDB
}

// NewUserHandler returns a new UserHandler backed by the given in-memory DB.
func NewUserHandler(db *memdb.MemDB) *UserHandler {
	return &UserHandler{db: db}
}

// Create creates a new user.
func (s *UserHandler) Create(_ context.Context, u *repository.User) error {
	txn := s.db.Txn(true)
	u.ID = uuid.New()
	record := &userRecord{
		ID:       u.ID.String(),
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
	}
	if err := txn.Insert("users", record); err != nil {
		txn.Abort()
		return fmt.Errorf("failed to insert user: %w", err)
	}
	txn.Commit()
	return nil
}

// FindAll returns all users.
func (s *UserHandler) FindAll(_ context.Context) ([]repository.User, error) {
	txn := s.db.Txn(false)
	it, err := txn.Get("users", "id")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	var users []repository.User
	for obj := it.Next(); obj != nil; obj = it.Next() {
		r := obj.(*userRecord)
		users = append(users, repository.User{
			ID:       uuid.MustParse(r.ID),
			Username: r.Username,
			Email:    r.Email,
			Name:     r.Name,
		})
	}
	return users, nil
}

// FindByID returns a user by ID.
func (s *UserHandler) FindByID(_ context.Context, id string) (*repository.User, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First("users", "id", id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if raw == nil {
		return nil, repository.ErrNotFound
	}
	r := raw.(*userRecord)
	return &repository.User{
		ID:       uuid.MustParse(r.ID),
		Username: r.Username,
		Email:    r.Email,
		Name:     r.Name,
	}, nil
}

// FindByUsername returns a user by username.
func (s *UserHandler) FindByUsername(_ context.Context, username string) (*repository.User, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First("users", "username", username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	if raw == nil {
		return nil, repository.ErrNotFound
	}
	r := raw.(*userRecord)
	return &repository.User{
		ID:       uuid.MustParse(r.ID),
		Username: r.Username,
		Email:    r.Email,
		Name:     r.Name,
	}, nil
}

// FindByEmail returns a user by email.
func (s *UserHandler) FindByEmail(_ context.Context, email string) (*repository.User, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First("users", "email", email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	if raw == nil {
		return nil, repository.ErrNotFound
	}
	r := raw.(*userRecord)
	return &repository.User{
		ID:       uuid.MustParse(r.ID),
		Username: r.Username,
		Email:    r.Email,
		Name:     r.Name,
	}, nil
}
