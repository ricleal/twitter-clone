package entities

import "errors"

var (
	// ErrInvalidTweet is returned when a user has an invalid email.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUserID is returned when a user has an invalid ID.
	ErrInvalidUserID = errors.New("invalid user id")
)
