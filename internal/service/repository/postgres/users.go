package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/orm"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserServer struct {
	dbConn repository.DBTx
}

func NewUserServer(dbConn repository.DBTx) *UserServer {
	return &UserServer{
		dbConn: dbConn,
	}
}

func (s *UserServer) Create(ctx context.Context, u *repository.User) error {
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

func (s *UserServer) FindAll(ctx context.Context) ([]repository.User, error) {
	ormUsers, err := orm.Users().All(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	var users []repository.User
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

func (s *UserServer) FindByID(ctx context.Context, id string) (*repository.User, error) {
	ormUser, err := orm.FindUser(ctx, s.dbConn, id)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
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

func (s *UserServer) FindByUsername(ctx context.Context, username string) (*repository.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.Username.EQ(username)).One(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
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

func (s *UserServer) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.Email.EQ(email)).One(ctx, s.dbConn)
	if err != nil {
		// Check if the error is a not found error
		if err.Error() == "sql: no rows in result set" {
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
