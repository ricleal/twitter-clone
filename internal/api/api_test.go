//go:build integration
// +build integration

package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/api"
	"github.com/ricleal/twitter-clone/internal/api/openapi"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/store"
	"github.com/ricleal/twitter-clone/testhelpers"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	server *httptest.Server
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (ts *APITestSuite) SetupTest() {

	// Set up our data store
	store := store.NewMemStore()
	st := service.NewTweetService(store)
	su := service.NewUserService(store)
	// set up our API
	twitterAPI := api.New(su, st)
	r := chi.NewRouter()
	openapi.HandlerFromMux(twitterAPI, r)
	ts.server = httptest.NewServer(r)
}

func (ts *APITestSuite) TearDownTest() {
	ts.server.Close()
}

func (ts *APITestSuite) TestGetUsersEmpty() {
	ctx := context.Background()
	var response struct{}

	statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
	ts.Require().NoError(err)
	ts.Require().Equal(http.StatusNoContent, statusCode)

}

func (ts *APITestSuite) TestCreateAndGetUser() {
	ctx := context.Background()

	userStr := `{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }`
	var userID string
	{
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/users", userStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	}
	{
		var response []map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Len(response, 1)
		ts.Require().Equal("foo", response[0]["username"])
		ts.Require().Equal("John Doe", response[0]["name"])
		ts.Require().Equal("jd@mail.com", response[0]["email"])
		userID = response[0]["id"].(string)
	}
	{
		var response map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users/"+userID, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		ts.Require().Equal("foo", response["username"])
		ts.Require().Equal("John Doe", response["name"])
		ts.Require().Equal("jd@mail.com", response["email"])
	}
	{
		// invalid user id
		var response struct{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users/"+uuid.NewString(), &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusNoContent, statusCode)
	}

}

func (ts *APITestSuite) TestCreateAndGetTweets() {
	ctx := context.Background()

	userStr := `{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }`
	var userID string
	{
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/users", userStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)
	}
	{
		var response []map[string]interface{}
		statusCode, err := testhelpers.Get(ctx, ts.server.URL+"/users", &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusOK, statusCode)
		userID = response[0]["id"].(string)
	}
	{
		tweetStr := `{ "user_id": "` + userID + `", "content": "Hello World" }`
		var response struct{}
		statusCode, err := testhelpers.Post(ctx, ts.server.URL+"/tweets", tweetStr, &response)
		ts.Require().NoError(err)
		ts.Require().Equal(http.StatusCreated, statusCode)

	}
}
