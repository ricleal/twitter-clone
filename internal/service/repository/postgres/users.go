package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/orm"
)

// UserStorage is a postgres implementation of the repository.UserStorage interface.
type UserStorage struct {
	dbConn repository.DBTx
}

// NewUserStorage returns a new UserServer.
func NewUserStorage(dbConn repository.DBTx) *UserStorage {
	return &UserStorage{
		dbConn: dbConn,
	}
}

// Create creates a new user.
func (s *UserStorage) Create(ctx context.Context, u *repository.User) error {
	user := orm.User{
		ID:       uuid.NewString(),
		Username: u.Username,
		Email:    u.Email,
		Name:     null.StringFrom(u.Name),
	}

	err := user.Insert(ctx, s.dbConn, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// FindAll returns all users.
func (s *UserStorage) FindAll(ctx context.Context) ([]repository.User, error) {
	ormUsers, err := orm.Users().All(ctx, s.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	users := make([]repository.User, 0, len(ormUsers))
	for _, u := range ormUsers {
		users = append(users, repository.User{
			ID:       uuid.MustParse(u.ID),
			Username: u.Username,
			Email:    u.Email,
			Name:     u.Name.String,
		})
	}

	return users, nil
}

// FindByID returns a user by ID.
func (s *UserStorage) FindByID(ctx context.Context, id string) (*repository.User, error) {
	ormUser, err := orm.FindUser(ctx, s.dbConn, id)
	if err != nil {
		// Check if the error is a not found error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.String,
	}, nil
}

// FindByUsername returns a user by username.
func (s *UserStorage) FindByUsername(ctx context.Context, username string) (*repository.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.Username.EQ(username)).One(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.String,
	}, nil
}

// FindByEmail returns a user by email.
func (s *UserStorage) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.Email.EQ(email)).One(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &repository.User{
		ID:       uuid.MustParse(ormUser.ID),
		Username: ormUser.Username,
		Email:    ormUser.Email,
		Name:     ormUser.Name.String,
	}, nil
}
