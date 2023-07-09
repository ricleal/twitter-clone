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

type TweetsTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	s         *postgres.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestTweetsTestSuite(t *testing.T) {
	suite.Run(t, new(TweetsTestSuite))
}

func (ts *TweetsTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = test.SetupDB(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = postgres.NewStorage(ctx)
	require.NoError(ts.T(), err)
}

func (ts *TweetsTestSuite) TearDownTest() {
	ctx := context.Background()
	err := test.TeardownDB(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *TweetsTestSuite) TestValid() {
	s := store.NewPersistentStore(ts.s.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)
	ctx := context.Background()

	ts.Run("get all users empty DB", func() {
		users, err := su.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(users, 0)
	})

	ts.Run("get an user by id empty DB", func() {
		user, err := su.FindByID(ctx, uuid.New().String())
		ts.Require().NoError(err)
		ts.Require().Nil(user)

		err = su.Create(ctx, &entities.User{
			Username: "test",
			Email:    "test@test.com",
			Name:     "John Doe",
		})
		ts.Require().NoError(err)
	})

	var userID uuid.UUID
	ts.Run("get all users", func() {
		users, err := su.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(users, 1)
		userID = users[0].ID
	})

	ts.Run("create a tweet", func() {
		err := st.Create(ctx, &entities.Tweet{
			UserID:  userID,
			Content: "Hello World",
		})
		ts.Require().NoError(err)
	})

	ts.Run("create a tweet with invalid user", func() {
		err := st.Create(ctx, &entities.Tweet{
			UserID:  uuid.New(),
			Content: "user does not exist",
		})
		ts.Require().Error(err)
		ts.Require().Contains(err.Error(), "invalid user id")
	})

	ts.Run("get all tweets", func() {
		tweets, err := st.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 1)
	})

	ts.Run("create a tweet with invalid user", func() {
		err := st.Create(ctx, &entities.Tweet{
			UserID:  uuid.New(),
			Content: "Hello World",
		})
		ts.Require().ErrorIs(err, entities.ErrInvalidUserID)
	})
}

func (ts *TweetsTestSuite) TestInvalid() {
	s := store.NewPersistentStore(ts.s.DB())
	st := service.NewTweetService(s)
	ctx := context.Background()

	// create a tweet with invalid user
	err := st.Create(ctx, &entities.Tweet{
		UserID:  uuid.New(),
		Content: "Hello World",
	})
	ts.Require().Error(err)
}
