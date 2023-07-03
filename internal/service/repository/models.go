package repository

import "github.com/google/uuid"

// User represents a user in the storage system.
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Name     string
}

// Tweet represents a tweet in the storage system.
type Tweet struct {
	ID      uuid.UUID
	Content string
	UserID  uuid.UUID
}
