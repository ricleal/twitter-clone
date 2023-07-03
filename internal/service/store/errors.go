package store

import "errors"

// ErrExecTxError is returned when ExecTx fails.
var ErrExecTxError = errors.New("ExecTxError")

// ExecTxError is returned when ExecTx fails.
type ExecTxError struct {
	message string
}

// Error returns the error message.
func (r *ExecTxError) Error() string {
	return r.message
}

// Unwrap returns the wrapped error.
func (r *ExecTxError) Unwrap() error {
	return ErrExecTxError
}

// NewExecTxError returns a ExecTxError with the given message.
func NewExecTxError(message string) *ExecTxError {
	return &ExecTxError{
		message: "ExecTxError: " + message,
	}
}
