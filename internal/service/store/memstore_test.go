package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

func TestMemStore_Tweets(t *testing.T) {
	memStore := store.NewMemStore()

	// Ensure the returned TweetRepository is the memory implementation
	tweetRepo := memStore.Tweets()
	_, ok := tweetRepo.(*memory.TweetHandler)
	if !ok {
		t.Error("Expected TweetRepository to be a memory implementation")
	}
}

func TestMemStore_Users(t *testing.T) {
	memStore := store.NewMemStore()

	// Ensure the returned UserRepository is the memory implementation
	userRepo := memStore.Users()
	_, ok := userRepo.(*memory.UserHandler)
	if !ok {
		t.Error("Expected UserRepository to be a memory implementation")
	}
}

func TestMemStore_ExecTx_Success(t *testing.T) {
	memStore := store.NewMemStore()

	// Test a successful transaction execution
	err := memStore.ExecTx(context.Background(), func(s store.Store) error {
		// Noop transaction
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestMemStore_ExecTx_Error(t *testing.T) {
	memStore := store.NewMemStore()
	memStore.TransactionError = true

	// Test an error during transaction execution
	err := memStore.ExecTx(context.Background(), func(s store.Store) error {
		// Noop transaction
		return nil
	})

	expectedErr := store.NewExecTxError("a transaction related error occurred")

	// Verify the expected error is returned
	if !errors.As(err, &expectedErr) {
		t.Errorf("Expected error '%v', got: '%v'", expectedErr, err)
	}
}
