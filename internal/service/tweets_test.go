//go:build integration
// +build integration

package service_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

type TeetsTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	ctx       context.Context
	s         *postgres.Handler
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestTeetsTestSuite(t *testing.T) {
	suite.Run(t, new(TeetsTestSuite))
}

func (ts *TeetsTestSuite) SetupTest() {
	var err error
	ts.ctx = context.Background()
	ts.container, err = test.SetupDB(ts.ctx)
	require.NoError(ts.T(), err)
	ts.s, err = postgres.NewHandler(ts.ctx)
	require.NoError(ts.T(), err)
}

func (ts *TeetsTestSuite) TearDownTest() {
	err := test.TeardownDB(ts.ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *TeetsTestSuite) TestValid() {
	s := store.NewSQLStore(ts.s.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)

	// get all users empty DB
	users, err := su.FindAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 0)

	// get an user by id empty DB
	user, err := su.FindByID(ts.ctx, uuid.New().String())
	ts.Require().NoError(err)
	ts.Require().Nil(user)

	err = su.Create(ts.ctx, &entities.User{
		Username: "test",
		Email:    "test@test.com",
		Name:     "John Doe",
	})
	ts.Require().NoError(err)
	// get all users
	users, err = su.FindAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 1)

	// create a tweet
	err = st.Create(ts.ctx, &entities.Tweet{
		UserID:  users[0].ID,
		Content: "Hello World",
	})
	ts.Require().NoError(err)

	// get all tweets
	tweets, err := st.FindAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tweets, 1)

	// create a tweet with invalid user
	err = st.Create(ts.ctx, &entities.Tweet{
		UserID:  uuid.New(),
		Content: "Hello World",
	})
	ts.Require().ErrorIs(err, entities.ErrInvalidUserID)
}

func (ts *TeetsTestSuite) TestInvalid() {
	s := store.NewSQLStore(ts.s.DB())
	st := service.NewTweetService(s)

	// create a tweet with invalid user
	err := st.Create(ts.ctx, &entities.Tweet{
		UserID:  uuid.New(),
		Content: "Hello World",
	})
	ts.Require().Error(err)
}
