package memory

import (
	"context"
	"fmt"

	memdb "github.com/hashicorp/go-memdb"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// TweetHandler is a memory implementation of the repository.TweetRepository interface.
type TweetHandler struct {
	db *memdb.MemDB
}

// NewTweetHandler returns a new TweetHandler backed by the given in-memory DB.
func NewTweetHandler(db *memdb.MemDB) *TweetHandler {
	return &TweetHandler{db: db}
}

// Create creates a new tweet.
func (s *TweetHandler) Create(_ context.Context, t *repository.Tweet) error {
	txn := s.db.Txn(true)
	t.ID = uuid.New()
	record := &tweetRecord{
		ID:      t.ID.String(),
		Content: t.Content,
		UserID:  t.UserID.String(),
	}
	if err := txn.Insert("tweets", record); err != nil {
		txn.Abort()
		return fmt.Errorf("failed to insert tweet: %w", err)
	}
	txn.Commit()
	return nil
}

// FindAll returns all tweets.
func (s *TweetHandler) FindAll(_ context.Context) ([]repository.Tweet, error) {
	txn := s.db.Txn(false)
	it, err := txn.Get("tweets", "id")
	if err != nil {
		return nil, fmt.Errorf("failed to get tweets: %w", err)
	}
	var tweets []repository.Tweet
	for obj := it.Next(); obj != nil; obj = it.Next() {
		r := obj.(*tweetRecord)
		tweets = append(tweets, repository.Tweet{
			ID:      uuid.MustParse(r.ID),
			Content: r.Content,
			UserID:  uuid.MustParse(r.UserID),
		})
	}
	return tweets, nil
}

// FindByID returns a tweet by ID.
func (s *TweetHandler) FindByID(_ context.Context, id string) (*repository.Tweet, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First("tweets", "id", id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tweet: %w", err)
	}
	if raw == nil {
		return nil, repository.ErrNotFound
	}
	r := raw.(*tweetRecord)
	return &repository.Tweet{
		ID:      uuid.MustParse(r.ID),
		Content: r.Content,
		UserID:  uuid.MustParse(r.UserID),
	}, nil
}
