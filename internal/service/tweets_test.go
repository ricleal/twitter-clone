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
	s         *postgres.Handler
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestTeetsTestSuite(t *testing.T) {
	suite.Run(t, new(TeetsTestSuite))
}

func (ts *TeetsTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = test.SetupDB(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = postgres.NewHandler(ctx)
	require.NoError(ts.T(), err)
}

func (ts *TeetsTestSuite) TearDownTest() {
	ctx := context.Background()
	err := test.TeardownDB(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *TeetsTestSuite) TestValid() {
	s := store.NewSQLStore(ts.s.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)
	ctx := context.Background()

	// get all users empty DB
	users, err := su.FindAll(ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 0)

	// get an user by id empty DB
	user, err := su.FindByID(ctx, uuid.New().String())
	ts.Require().NoError(err)
	ts.Require().Nil(user)

	err = su.Create(ctx, &entities.User{
		Username: "test",
		Email:    "test@test.com",
		Name:     "John Doe",
	})
	ts.Require().NoError(err)
	// get all users
	users, err = su.FindAll(ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 1)

	// create a tweet
	err = st.Create(ctx, &entities.Tweet{
		UserID:  users[0].ID,
		Content: "Hello World",
	})
	ts.Require().NoError(err)

	// get all tweets
	tweets, err := st.FindAll(ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tweets, 1)

	// create a tweet with invalid user
	err = st.Create(ctx, &entities.Tweet{
		UserID:  uuid.New(),
		Content: "Hello World",
	})
	ts.Require().ErrorIs(err, entities.ErrInvalidUserID)
}

func (ts *TeetsTestSuite) TestInvalid() {
	s := store.NewSQLStore(ts.s.DB())
	st := service.NewTweetService(s)
	ctx := context.Background()

	// create a tweet with invalid user
	err := st.Create(ctx, &entities.Tweet{
		UserID:  uuid.New(),
		Content: "Hello World",
	})
	ts.Require().Error(err)
}
