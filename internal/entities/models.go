package entities

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Name     string
}

type Tweet struct {
	ID      uuid.UUID
	Content string
	UserID  uuid.UUID
}
