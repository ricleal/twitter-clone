package repository

import "errors"

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("not found")
