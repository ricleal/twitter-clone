package entities

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalidEmail is returned when a user has an invalid email according to the regex.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUserID is returned when creating a tweet a user has an invalid ID.
	ErrInvalidUserID = errors.New("invalid user id")
)
