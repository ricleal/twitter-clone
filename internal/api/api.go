package api

import (
	"encoding/json"
	"errors"
	"net/http"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/api/openapi"
	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service"
)

// 200 – OK
// Everything is working, The resource has been fetched and is transmitted in the message body.

// 201 – CREATED
// A new resource has been created

// 204 – NO CONTENT
// The resource was successfully deleted, no response body

// 304 – NOT MODIFIED
// This is used for caching purposes. It tells the client that the response has not been modified, so the client can continue to use the same cached version of the response.

// 400 – BAD REQUEST
// The request was invalid or cannot be served. The exact error should be explained in the error payload.

// 401 – UNAUTHORIZED
// The request requires user authentication.

// 403 – FORBIDDEN
// The server understood the request but is refusing it or the access is not allowed.

// 404 – NOT FOUND
// There is no resource behind the URI.

// 500 – INTERNAL SERVER ERROR API
// If an error occurs in the global catch blog, the stack trace should be logged and not returned as a response.

type twitterServer struct {
	tweetService service.TweetService
	userService  service.UserService
}

func New(userService service.UserService, tweetService service.TweetService) *twitterServer {
	return &twitterServer{
		tweetService: tweetService,
		userService:  userService,
	}
}

// List all tweets
// (GET /tweets)
func (t *twitterServer) GetTweets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tweets, err := t.tweetService.FindAll(ctx)
	if err != nil {
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error listing tweets", err)
		return
	}
	if len(tweets) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var openapiTweets []openapi.Tweet
	for _, tweet := range tweets {
		openapiTweets = append(openapiTweets, openapi.Tweet{
			Id:      &tweet.ID,
			Content: tweet.Content,
			UserId:  tweet.UserID,
		})
	}

	json.NewEncoder(w).Encode(openapiTweets)
}

// Create a tweet
// (POST /tweets)
func (t *twitterServer) PostTweets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newTweet openapi.Tweet
	if err := json.NewDecoder(r.Body).Decode(&newTweet); err != nil {
		sendAPIError(ctx, w, http.StatusBadRequest, "Invalid format for Tweet", err)
		return
	}

	// convert openapi.Tweet to entity.Tweet
	tweet := &entities.Tweet{
		Content: newTweet.Content,
		UserID:  newTweet.UserId,
	}

	if err := t.tweetService.Create(ctx, tweet); err != nil {
		if errors.Is(err, entities.ErrInvalidUserID) {
			sendAPIError(ctx, w, http.StatusNoContent, "Invalid user ID", err)
			return
		}
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error creating tweet", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get tweet by ID
// (GET /tweets/{id})
func (t *twitterServer) GetTweetsId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx := r.Context()

	tweet, err := t.tweetService.FindByID(ctx, id.String())
	if err != nil {
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error getting tweet", err)
		return
	}
	if tweet == nil {
		sendAPIError(ctx, w, http.StatusNoContent, "Tweet not found", nil)
		return
	}

	openapiTweet := openapi.Tweet{
		Id:      &tweet.ID,
		Content: tweet.Content,
		UserId:  tweet.UserID,
	}

	json.NewEncoder(w).Encode(openapiTweet)
}

// List all users
// (GET /users)
func (t *twitterServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := t.userService.FindAll(ctx)
	if err != nil {
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error listing users", err)
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// convert entity.User to openapi.User
	var openapiUsers []*openapi.User
	for _, user := range users {
		openapiUsers = append(openapiUsers, &openapi.User{
			Id:       &user.ID,
			Username: user.Username,
			Email:    openapi_types.Email(user.Email),
		})
		if user.Name != "" {
			openapiUsers[len(openapiUsers)-1].Name = &user.Name
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(openapiUsers)

}

// Create a user
// (POST /users)
func (t *twitterServer) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newUser openapi.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendAPIError(ctx, w, http.StatusBadRequest, "Invalid format for User", err)
		return
	}

	// convert openapi.User to entity.User
	user := &entities.User{
		Username: newUser.Username,
		Email:    string(newUser.Email),
	}
	if newUser.Name != nil {
		user.Name = *newUser.Name
	}

	if err := t.userService.Create(ctx, user); err != nil {
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get user profile by ID
// (GET /users/{id})
func (t *twitterServer) GetUsersId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	ctx := r.Context()

	user, err := t.userService.FindByID(ctx, id.String())
	if err != nil {
		sendAPIError(ctx, w, http.StatusInternalServerError, "Error getting user", err)
		return
	}
	if user == nil {
		sendAPIError(ctx, w, http.StatusNoContent, "User not found", nil)
		return
	}

	openapiUser := openapi.User{
		Id:       &user.ID,
		Username: user.Username,
		Email:    openapi_types.Email(user.Email),
	}
	if user.Name != "" {
		openapiUser.Name = &user.Name
	}

	json.NewEncoder(w).Encode(openapiUser)
}
