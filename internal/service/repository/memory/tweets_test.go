//go:build !integration
// +build !integration

package memory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

func TestTweetHandlerCreate(t *testing.T) {
	tweetHandler := memory.NewTweetHandler()

	// Test creating a new tweet
	tweet := &repository.Tweet{
		Content: "Hello, world!",
	}

	err := tweetHandler.Create(context.Background(), tweet)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Verify the tweet has a valid ID
	if tweet.ID == uuid.Nil {
		t.Error("Invalid tweet ID")
	}

	// Verify the tweet is stored in the handler
	tweets, err := tweetHandler.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error retrieving tweets: %v", err)
	}

	if len(tweets) != 1 {
		t.Errorf("Expected 1 tweet, got %d", len(tweets))
	}

	if tweets[0].Content != tweet.Content {
		t.Errorf("Expected tweet content %q, got %q", tweet.Content, tweets[0].Content)
	}
}

func TestTweetHandlerFindAll(t *testing.T) {
	tweetHandler := memory.NewTweetHandler()

	// Create some tweets
	tweet1 := &repository.Tweet{
		Content: "Tweet 1",
	}

	tweet2 := &repository.Tweet{
		Content: "Tweet 2",
	}

	err := tweetHandler.Create(context.Background(), tweet1)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	err = tweetHandler.Create(context.Background(), tweet2)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Retrieve all tweets
	tweets, err := tweetHandler.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error retrieving tweets: %v", err)
	}

	// Verify the number of tweets
	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets, got %d", len(tweets))
	}

	// Verify the tweet contents
	if tweets[0].Content != tweet1.Content {
		t.Errorf("Expected tweet content %q, got %q", tweet1.Content, tweets[0].Content)
	}

	if tweets[1].Content != tweet2.Content {
		t.Errorf("Expected tweet content %q, got %q", tweet2.Content, tweets[1].Content)
	}
}

func TestTweetHandlerFindByID(t *testing.T) {
	tweetHandler := memory.NewTweetHandler()

	// Create a tweet
	tweet := &repository.Tweet{
		Content: "Hello, world!",
	}

	err := tweetHandler.Create(context.Background(), tweet)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Retrieve the tweet by ID
	foundTweet, err := tweetHandler.FindByID(context.Background(), tweet.ID.String())
	if err != nil {
		t.Errorf("Error retrieving tweet: %v", err)
	}

	// Verify the retrieved tweet is the same as the original tweet
	if foundTweet.Content != tweet.Content {
		t.Errorf("Expected tweet content %q, got %q", tweet.Content, foundTweet.Content)
	}
}

func TestTweetHandlerFindByIDNotFound(t *testing.T) {
	tweetHandler := memory.NewTweetHandler()

	// Retrieve a non-existent tweet by ID
	_, err := tweetHandler.FindByID(context.Background(), "non-existent-id")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
