//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
)

type TweetsTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	s         *postgres.Handler
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
	ts.s, err = postgres.NewHandler(ctx)
	require.NoError(ts.T(), err)
}

func (ts *TweetsTestSuite) TearDownTest() {
	ctx := context.Background()
	err := test.TeardownDB(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *TweetsTestSuite) TestData() {
	ctx := context.Background()
	t := postgres.NewTweetServer(ts.s.DB())
	u := postgres.NewUserServer(ts.s.DB())

	// Find all tweets empty DB
	{
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 0)
	}
	// Find tweet empty DB
	{
		tweet, err := t.FindByID(ctx, uuid.New().String())
		ts.Require().ErrorIs(err, repository.ErrNotFound)
		ts.Require().Nil(tweet)
	}

	{
		err := u.Create(ctx, &repository.User{
			Username: "test",
			Email:    "test@test.com",
		})
		ts.Require().NoError(err)
	}
	// Find all users
	{
		users, err := u.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(users, 1)
	}

	var userID string
	// Find user by Username
	{
		user, err := u.FindByUsername(ctx, "test")
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)
		userID = user.ID.String()
	}
	// Find user by ID
	{
		user, err := u.FindByID(ctx, userID)
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)
	}

	// Find all tweets with 0 tweet
	{
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 0)
	}
	// create a tweet
	{
		tweet := &repository.Tweet{
			UserID:  uuid.Must(uuid.Parse(userID)),
			Content: "Lorem ipsum dolor sit amet",
		}
		err := t.Create(ctx, tweet)
		ts.Require().NoError(err)
	}
	// Find all tweets with 1 tweet
	{
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 1)
	}
	// create a tweet
	{
		tweet2 := &repository.Tweet{
			UserID:  uuid.Must(uuid.Parse(userID)),
			Content: "Ut enim ad minim veniam",
		}
		err := t.Create(ctx, tweet2)
		ts.Require().NoError(err)
	}
	var tweetID string
	// Find all tweets
	{
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 2)
		tweetID = tweets[0].ID.String()
	}
	// Find tweet by ID
	{
		tweet, err := t.FindByID(ctx, tweetID)
		ts.Require().NoError(err)
		ts.Require().NotNil(tweet)
		ts.Require().Equal(36, len(tweet.ID.String()))
	}
	// Find invalid tweet by ID"
	{
		tweet, err := t.FindByID(ctx, uuid.New().String())
		ts.Require().ErrorIs(err, repository.ErrNotFound)
		ts.Require().Nil(tweet)
	}
}
