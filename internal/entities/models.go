package entities

import "github.com/google/uuid"

// User represents a user in the domain.
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Name     string
}

// Tweet represents a tweet in the domain.
type Tweet struct {
	ID      uuid.UUID
	Content string
	UserID  uuid.UUID
}
