package memory_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

func newTestTweetHandler(t *testing.T) *memory.TweetHandler {
	t.Helper()
	db, err := memory.NewDB()
	if err != nil {
		t.Fatalf("failed to create in-memory DB: %v", err)
	}
	return memory.NewTweetHandler(db)
}

func TestTweetHandlerCreate(t *testing.T) {
	tweetHandler := newTestTweetHandler(t)

	// Test creating a new tweet
	tweet := &entities.Tweet{
		Content: "Hello, world!",
	}

	err := tweetHandler.Create(t.Context(), tweet)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Verify the tweet has a valid ID
	if tweet.ID == uuid.Nil {
		t.Error("Invalid tweet ID")
	}

	// Verify the tweet is stored in the handler
	tweets, err := tweetHandler.FindAll(t.Context())
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
	tweetHandler := newTestTweetHandler(t)

	// Create some tweets
	tweet1 := &entities.Tweet{
		Content: "Tweet 1",
	}

	tweet2 := &entities.Tweet{
		Content: "Tweet 2",
	}

	err := tweetHandler.Create(t.Context(), tweet1)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	err = tweetHandler.Create(t.Context(), tweet2)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Retrieve all tweets
	tweets, err := tweetHandler.FindAll(t.Context())
	if err != nil {
		t.Errorf("Error retrieving tweets: %v", err)
	}

	// Verify the number of tweets
	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets, got %d", len(tweets))
	}

	// Verify both tweet contents are present (order not guaranteed)
	contentSet := make(map[string]bool)
	for _, tw := range tweets {
		contentSet[tw.Content] = true
	}
	if !contentSet[tweet1.Content] {
		t.Errorf("Expected tweet content %q to be present", tweet1.Content)
	}
	if !contentSet[tweet2.Content] {
		t.Errorf("Expected tweet content %q to be present", tweet2.Content)
	}
}

func TestTweetHandlerFindByID(t *testing.T) {
	tweetHandler := newTestTweetHandler(t)

	// Create a tweet
	tweet := &entities.Tweet{
		Content: "Hello, world!",
	}

	err := tweetHandler.Create(t.Context(), tweet)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}

	// Retrieve the tweet by ID
	foundTweet, err := tweetHandler.FindByID(t.Context(), tweet.ID.String())
	if err != nil {
		t.Errorf("Error retrieving tweet: %v", err)
	}

	// Verify the retrieved tweet is the same as the original tweet
	if foundTweet.Content != tweet.Content {
		t.Errorf("Expected tweet content %q, got %q", tweet.Content, foundTweet.Content)
	}
}

func TestTweetHandlerFindByIDNotFound(t *testing.T) {
	tweetHandler := newTestTweetHandler(t)

	// Retrieve a non-existent tweet by ID
	_, err := tweetHandler.FindByID(t.Context(), "non-existent-id")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
