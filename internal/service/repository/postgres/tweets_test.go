//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TweetsTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	ctx       context.Context
	s         *postgres.Server
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestTweetsTestSuite(t *testing.T) {
	suite.Run(t, new(TweetsTestSuite))
}

func (ts *TweetsTestSuite) SetupTest() {
	var err error
	ts.ctx = context.Background()
	ts.container, err = test.SetupDB(ts.ctx)
	require.NoError(ts.T(), err)
	ts.s = postgres.New()
}

func (ts *TweetsTestSuite) TearDownTest() {
	test.TeardownDB(ts.ctx, ts.container)
	ts.s.Close()
}

func (ts *TweetsTestSuite) TestData() {

	t := postgres.NewTweetServer(ts.s.DB())

	u := postgres.NewUserServer(ts.s.DB())

	// create a user
	err := u.Create(ts.ctx, &repository.User{
		Username: "test",
		Email:    "test@test.com",
	})
	ts.Require().NoError(err)

	// Find all users
	users, err := u.FindAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(users, 1)

	// Find user by Username
	user, err := u.FindByUsername(ts.ctx, "test")
	ts.Require().NoError(err)
	ts.Require().Equal("test", user.Username)

	// Find user by ID
	user, err = u.FindByID(ts.ctx, user.ID.String())
	ts.Require().NoError(err)
	ts.Require().Equal("test", user.Username)

	// create a tweet
	tweet := &repository.Tweet{
		UserID:  user.ID,
		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}
	err = t.Create(ts.ctx, tweet)
	ts.Require().NoError(err)

	// create a tweet
	tweet2 := &repository.Tweet{
		UserID:  user.ID,
		Content: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	}
	err = t.Create(ts.ctx, tweet2)
	ts.Require().NoError(err)

	// Find all tweets
	tweets, err := t.FindAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(tweets, 2)

	// Find tweet by ID
	tweet, err = t.FindByID(ts.ctx, tweets[0].ID.String())
	ts.Require().NoError(err)
	ts.Require().Equal("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", tweet.Content)

}
