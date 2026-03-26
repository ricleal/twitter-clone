package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/models"
)

// UserStorage is a postgres implementation of the repository.UserStorage interface.
type UserStorage struct {
	dbConn bob.Executor
}

// NewUserStorage returns a new UserServer.
func NewUserStorage(dbConn bob.Executor) *UserStorage {
	return &UserStorage{
		dbConn: dbConn,
	}
}

// Create creates a new user.
func (s *UserStorage) Create(ctx context.Context, u *repository.User) error {
	setter := &models.UserSetter{
		ID:       omit.From(uuid.NewString()),
		Username: omit.From(u.Username),
		Email:    omit.From(u.Email),
		Name:     omitnull.From(u.Name),
	}

	_, err := models.Users.Insert(setter).One(ctx, s.dbConn)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// FindAll returns all users.
func (s *UserStorage) FindAll(ctx context.Context) ([]repository.User, error) {
	ormUsers, err := models.Users.Query().All(ctx, s.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	users := make([]repository.User, 0, len(ormUsers))
	for _, u := range ormUsers {
		users = append(users, repository.User{
			ID:       uuid.MustParse(u.ID),
			Username: u.Username,
			Email:    u.Email,
			Name:     u.Name.GetOrZero(),
		})
	}

	return users, nil
}

// FindByID returns a user by ID.
func (s *UserStorage) FindByID(ctx context.Context, id string) (*repository.User, error) {
	ormUser, err := models.FindUser(ctx, s.dbConn, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.GetOrZero(),
	}, nil
}

// FindByUsername returns a user by username.
func (s *UserStorage) FindByUsername(ctx context.Context, username string) (*repository.User, error) {
	ormUser, err := models.Users.Query(
		sm.Where(models.Users.Columns.Username.EQ(psql.Arg(username))),
	).One(ctx, s.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.GetOrZero(),
	}, nil
}

// FindByEmail returns a user by email.
func (s *UserStorage) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	ormUser, err := models.Users.Query(
		sm.Where(models.Users.Columns.Email.EQ(psql.Arg(email))),
	).One(ctx, s.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.GetOrZero(),
	}, nil
}
