package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/service"
)

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
		w.WriteHeader(http.StatusInternalServerError)
	}

	for _, tweet := range tweets {
		_ = tweet
	}
	w.Write([]byte("Hello world"))
}

// Create a tweet
// (POST /tweets)
func (t *twitterServer) PostTweets(w http.ResponseWriter, r *http.Request) {}

// Get tweet by ID
// (GET /tweets/{id})
func (t *twitterServer) GetTweetsId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {}

// List all users
// (GET /users)
func (t *twitterServer) GetUsers(w http.ResponseWriter, r *http.Request) {}

// Create a user
// (POST /users)
func (t *twitterServer) PostUsers(w http.ResponseWriter, r *http.Request) {}

// Get user profile by ID
// (GET /users/{id})
func (t *twitterServer) GetUsersId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {}
