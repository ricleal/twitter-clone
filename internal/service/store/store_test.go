//go:build integration
// +build integration

package store_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
	"github.com/ricleal/twitter-clone/internal/service/store"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type StoreTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	ctx       context.Context
	s         *postgres.Server
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (ts *StoreTestSuite) SetupTest() {
	var err error
	ts.ctx = context.Background()
	ts.container, err = test.SetupDB(ts.ctx)
	require.NoError(ts.T(), err)
	ts.s = postgres.New()
}

func (ts *StoreTestSuite) TearDownTest() {
	test.TeardownDB(ts.ctx, ts.container)
	ts.s.Close()
}

func (ts *StoreTestSuite) TestData() {

	s := store.NewSQLStore(ts.s.DB())

	if err := s.ExecTx(ts.ctx, func(s store.Store) error {

		tweetsRepo := s.Tweets()
		usersRepo := s.Users()

		// create a user
		err := usersRepo.Create(ts.ctx, &repository.User{
			Username: "test",
			Email:    "test@test.com",
		})
		ts.Require().NoError(err)

		// Find user by Username
		user, err := usersRepo.FindByUsername(ts.ctx, "test")
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)

		// create a tweet
		tweet := &repository.Tweet{
			UserID:  user.ID,
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		}
		err = tweetsRepo.Create(ts.ctx, tweet)
		ts.Require().NoError(err)

		// Find all tweets
		tweets, err := tweetsRepo.FindAll(ts.ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 1)
		return err
	}); err != nil {
		ts.Require().NoError(err)
	}

}
