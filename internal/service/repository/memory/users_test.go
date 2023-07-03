package memory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

func TestUserHandlerCreate(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Test creating a new user
	user := &repository.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Verify the user has a valid ID
	if user.ID == uuid.Nil {
		t.Error("Invalid user ID")
	}

	// Verify the user is stored in the handler
	users, err := userHandler.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error retrieving users: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0].Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, users[0].Username)
	}
}

func TestUserHandlerFindAll(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Create some users
	user1 := &repository.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	user2 := &repository.User{
		Username: "jane_doe",
		Email:    "jane.doe@example.com",
	}

	err := userHandler.Create(context.Background(), user1)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	err = userHandler.Create(context.Background(), user2)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve all users
	users, err := userHandler.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error retrieving users: %v", err)
	}

	// Verify the number of users
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Verify the usernames
	if users[0].Username != user1.Username {
		t.Errorf("Expected username %q, got %q", user1.Username, users[0].Username)
	}

	if users[1].Username != user2.Username {
		t.Errorf("Expected username %q, got %q", user2.Username, users[1].Username)
	}
}

func TestUserHandlerFindByID(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Create a user
	user := &repository.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by ID
	foundUser, err := userHandler.FindByID(context.Background(), user.ID.String())
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, foundUser.Username)
	}
}

func TestUserHandlerFindByIDNotFound(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Retrieve a non-existent user by ID
	_, err := userHandler.FindByID(context.Background(), "non-existent-id")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestUserHandlerFindByUsername(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Create a user
	user := &repository.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by username
	foundUser, err := userHandler.FindByUsername(context.Background(), user.Username)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, foundUser.Username)
	}
}

func TestUserHandlerFindByUsernameNotFound(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Retrieve a non-existent user by username
	_, err := userHandler.FindByUsername(context.Background(), "non-existent-username")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestUserHandlerFindByEmail(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Create a user
	user := &repository.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by email
	foundUser, err := userHandler.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Email != user.Email {
		t.Errorf("Expected email %q, got %q", user.Email, foundUser.Email)
	}
}

func TestUserHandlerFindByEmailNotFound(t *testing.T) {
	userHandler := memory.NewUserHandler()

	// Retrieve a non-existent user by email
	_, err := userHandler.FindByEmail(context.Background(), "non-existent-email@example.com")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
