package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ricleal/twitter-clone/internal/openapi"
)

type TwitterServer struct {
}

func (t TwitterServer) GetTweets(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

// Create a tweet
// (POST /tweets)
func (t TwitterServer) PostTweets(w http.ResponseWriter, r *http.Request) {}

// Get tweet by ID
// (GET /tweets/{id})
func (t TwitterServer) GetTweetsId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {}

// List all users
// (GET /users)
func (t TwitterServer) GetUsers(w http.ResponseWriter, r *http.Request) {}

// Create a user
// (POST /users)
func (t TwitterServer) PostUsers(w http.ResponseWriter, r *http.Request) {}

// Get user profile by ID
// (GET /users/{id})
func (t TwitterServer) GetUsersId(w http.ResponseWriter, r *http.Request, id uuid.UUID) {}

func main() {
	s := TwitterServer{}
	h := openapi.Handler(s)

	http.ListenAndServe(":3000", h)

}
