//go:build integration
// +build integration

package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	api "github.com/ricleal/twitter-clone/internal/api/v1"
	"github.com/ricleal/twitter-clone/internal/api/v1/openapi"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
	"github.com/ricleal/twitter-clone/internal/service/store"
	"github.com/ricleal/twitter-clone/testhelpers"
)

type APITestIntegrationSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	s         *postgres.Storage
	server    *httptest.Server
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAPITestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(APITestIntegrationSuite))
}

func (ts *APITestIntegrationSuite) SetupTest() {
	var err error
	ctx := context.Background()
	ts.container, err = test.SetupDB(ctx)
	require.NoError(ts.T(), err)
	ts.s, err = postgres.NewStorage(ctx)
	require.NoError(ts.T(), err)

	s := store.NewPersistentStore(ts.s.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)

	// set up our API
	twitterAPI := api.New(su, st)
	r := chi.NewRouter()
	openapi.HandlerFromMux(twitterAPI, r)
	ts.server = httptest.NewServer(r)
}

func (ts *APITestIntegrationSuite) TearDownTest() {
	ctx := context.Background()
	err := test.TeardownDB(ctx, ts.container)
	require.NoError(ts.T(), err)
	ts.s.Close()
}

func (ts *APITestIntegrationSuite) TestGetUsersEmpty() {
	ctx := context.Background()
	var response []interface{}
	statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
	ts.Require().NoError(err)
	ts.Require().Equal(http.StatusOK, statusCode)
	ts.Require().Equal([]interface{}{}, response)
}

func (ts *APITestIntegrationSuite) TestCreateAndGetUser() {
	ctx := context.Background()

	userStr := `{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }`
	var userID string

	ts.Run("Get users empty DB", func() {
		var response []interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Equal([]interface{}{}, response)
	})
	ts.Run("Create user", func() {
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/users", userStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	})
	ts.Run("Get users", func() {
		var response []map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Len(response, 1)
		ts.Require().Equal("foo", response[0]["username"])
		ts.Require().Equal("John Doe", response[0]["name"])
		ts.Require().Equal("jd@mail.com", response[0]["email"])
		userID = response[0]["id"].(string) //nolint:errcheck,forcetypeassert  // this is a test
	})
	ts.Run("Get user", func() {
		var response map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users/"+userID, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Equal("foo", response["username"])
		ts.Require().Equal("John Doe", response["name"])
		ts.Require().Equal("jd@mail.com", response["email"])
	})
	ts.Run("Get invalid user", func() {
		var response interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users/"+uuid.NewString(), &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Nil(response)
	})
}

func (ts *APITestIntegrationSuite) TestCreateAndGetTweets() {
	ctx := context.Background()

	userStr := `{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }`
	var userID string

	ts.Run("Create user", func() {
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/users", userStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	})
	ts.Run("Get users", func() {
		var response []map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		userID = response[0]["id"].(string) //nolint:errcheck,forcetypeassert  // this is a test
	})
	ts.Run("Create tweet", func() {
		tweetStr := `{ "user_id": "` + userID + `", "content": "Hello World" }`
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/tweets", tweetStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	})
	ts.Run("Create tweet invalid user", func() {
		// invalid user id
		tweetStr := `{ "user_id": "` + uuid.NewString() + `", "content": "Hello World" }`
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/tweets", tweetStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusNoContent, statusCode)
	})
	var tweetID string
	ts.Run("Get tweets", func() {
		var response []openapi.Tweet
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/tweets", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Len(response, 1)
		ts.Require().Equal("Hello World", response[0].Content)
		tweetID = response[0].Id.String()
	})
	ts.Run("Create tweet 2", func() {
		tweetStr := `{ "user_id": "` + userID + `", "content": "Hello World 2" }`
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/tweets", tweetStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	})
	ts.Run("Get tweets", func() {
		var response []openapi.Tweet
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/tweets", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Len(response, 2)
		ts.Require().Equal("Hello World", response[0].Content)
		ts.Require().Equal("Hello World 2", response[1].Content)
	})
	ts.Run("Get tweet", func() {
		var response openapi.Tweet
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/tweets/"+tweetID, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Equal("Hello World", response.Content)
	})
	ts.Run("Get invalid tweet", func() {
		var response interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/tweets/"+uuid.NewString(), &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Nil(response)
	})
}
