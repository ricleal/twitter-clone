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

func (ts *TweetsTestSuite) TestData() {
	ctx := context.Background()
	t := postgres.NewTweetStorage(ts.s.DB())
	u := postgres.NewUserStorage(ts.s.DB())

	ts.Run("Find all tweets empty DB", func() {
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 0)
	})
	ts.Run("Find tweet empty DB", func() {
		tweet, err := t.FindByID(ctx, uuid.New().String())
		ts.Require().ErrorIs(err, repository.ErrNotFound)
		ts.Require().Nil(tweet)
	})
	ts.Run("Create user", func() {
		err := u.Create(ctx, &repository.User{
			Username: "test",
			Email:    "test@test.com",
		})
		ts.Require().NoError(err)
	})
	ts.Run("Find all users", func() {
		users, err := u.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(users, 1)
	})
	var userID string
	ts.Run("Find user by username", func() {
		user, err := u.FindByUsername(ctx, "test")
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)
		userID = user.ID.String()
	})
	ts.Run("Find user by ID", func() {
		user, err := u.FindByID(ctx, userID)
		ts.Require().NoError(err)
		ts.Require().Equal("test", user.Username)
	})
	ts.Run("Create tweet", func() {
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 0)
	})
	ts.Run("Create tweet", func() {
		tweet := &repository.Tweet{
			UserID:  uuid.Must(uuid.Parse(userID)),
			Content: "Lorem ipsum dolor sit amet",
		}
		err := t.Create(ctx, tweet)
		ts.Require().NoError(err)
	})
	ts.Run("Create tweet with non existing user", func() {
		tweet := &repository.Tweet{
			UserID:  uuid.New(),
			Content: "user does not exist",
		}
		err := t.Create(ctx, tweet)
		ts.Require().Contains(err.Error(), "violates foreign key constraint")
	})
	ts.Run("Find all tweets with 1 tweet", func() {
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 1)
	})
	ts.Run("Create tweet 2", func() {
		tweet2 := &repository.Tweet{
			UserID:  uuid.Must(uuid.Parse(userID)),
			Content: "Ut enim ad minim veniam",
		}
		err := t.Create(ctx, tweet2)
		ts.Require().NoError(err)
	})
	var tweetID string
	ts.Run("Find all tweets", func() {
		tweets, err := t.FindAll(ctx)
		ts.Require().NoError(err)
		ts.Require().Len(tweets, 2)
		tweetID = tweets[0].ID.String()
	})
	ts.Run("Find tweet by ID", func() {
		tweet, err := t.FindByID(ctx, tweetID)
		ts.Require().NoError(err)
		ts.Require().NotNil(tweet)
		ts.Require().Equal(36, len(tweet.ID.String()))
	})
	ts.Run("Find invalid tweet by ID", func() {
		tweet, err := t.FindByID(ctx, uuid.New().String())
		ts.Require().ErrorIs(err, repository.ErrNotFound)
		ts.Require().Nil(tweet)
	})
}
