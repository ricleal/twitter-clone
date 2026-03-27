package memory_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ricleal/twitter-clone/internal/entities"
	"github.com/ricleal/twitter-clone/internal/service/repository"
	"github.com/ricleal/twitter-clone/internal/service/repository/memory"
)

func newTestUserHandler(t *testing.T) *memory.UserHandler {
	t.Helper()
	db, err := memory.NewDB()
	if err != nil {
		t.Fatalf("failed to create in-memory DB: %v", err)
	}
	return memory.NewUserHandler(db)
}

func TestUserHandlerCreate(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Test creating a new user
	user := &entities.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(t.Context(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Verify the user has a valid ID
	if user.ID == uuid.Nil {
		t.Error("Invalid user ID")
	}

	// Verify the user is stored in the handler
	users, err := userHandler.FindAll(t.Context())
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
	userHandler := newTestUserHandler(t)

	// Create some users
	user1 := &entities.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	user2 := &entities.User{
		Username: "jane_doe",
		Email:    "jane.doe@example.com",
	}

	err := userHandler.Create(t.Context(), user1)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	err = userHandler.Create(t.Context(), user2)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve all users
	users, err := userHandler.FindAll(t.Context())
	if err != nil {
		t.Errorf("Error retrieving users: %v", err)
	}

	// Verify the number of users
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Verify both usernames are present (order not guaranteed)
	usernameSet := make(map[string]bool)
	for _, u := range users {
		usernameSet[u.Username] = true
	}
	if !usernameSet[user1.Username] {
		t.Errorf("Expected username %q to be present", user1.Username)
	}
	if !usernameSet[user2.Username] {
		t.Errorf("Expected username %q to be present", user2.Username)
	}
}

func TestUserHandlerFindByID(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Create a user
	user := &entities.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(t.Context(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by ID
	foundUser, err := userHandler.FindByID(t.Context(), user.ID.String())
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, foundUser.Username)
	}
}

func TestUserHandlerFindByIDNotFound(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Retrieve a non-existent user by ID
	_, err := userHandler.FindByID(t.Context(), "non-existent-id")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestUserHandlerFindByUsername(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Create a user
	user := &entities.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(t.Context(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by username
	foundUser, err := userHandler.FindByUsername(t.Context(), user.Username)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Username != user.Username {
		t.Errorf("Expected username %q, got %q", user.Username, foundUser.Username)
	}
}

func TestUserHandlerFindByUsernameNotFound(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Retrieve a non-existent user by username
	_, err := userHandler.FindByUsername(t.Context(), "non-existent-username")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestUserHandlerFindByEmail(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Create a user
	user := &entities.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	err := userHandler.Create(t.Context(), user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by email
	foundUser, err := userHandler.FindByEmail(t.Context(), user.Email)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	// Verify the retrieved user is the same as the original user
	if foundUser.Email != user.Email {
		t.Errorf("Expected email %q, got %q", user.Email, foundUser.Email)
	}
}

func TestUserHandlerFindByEmailNotFound(t *testing.T) {
	userHandler := newTestUserHandler(t)

	// Retrieve a non-existent user by email
	_, err := userHandler.FindByEmail(t.Context(), "non-existent-email@example.com")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
