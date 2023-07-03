//go:build integration
// +build integration

package store_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

type StoreTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	s         *postgres.Storage
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (ts *StoreTestSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = test.SetupDB(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = postgres.NewStorage(ctx)
	require.NoError(ts.T(), err)
}

func (ts *StoreTestSuite) TearDownTest() {
	ctx := context.Background()
	err := test.TeardownDB(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *StoreTestSuite) TestTransaction() {
	ctx := context.Background()
	s := store.NewPersistentStore(ts.s.DB())

	if err := s.ExecTx(ctx, func(s store.Store) error {
		tweetsRepo := s.Tweets()
		usersRepo := s.Users()

		// create a user
		err := usersRepo.Create(ctx, &repository.User{
			Username: "test",
			Email:    "test@test.com",
		})
		ts.Require().NoError(err)

		// Find user by Username
		user, err := usersRepo.FindByUsername(ctx, "test")
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)

		// create a tweet
		tweet := &repository.Tweet{
			UserID:  user.ID,
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		}
		err = tweetsRepo.Create(ctx, tweet)
		ts.Require().NoError(err)

		// Find all tweets
		tweets, err := tweetsRepo.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 1)
		return nil
	}); err != nil {
		ts.Require().NoError(err)
	}
}
