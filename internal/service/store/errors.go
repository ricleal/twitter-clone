package store

import "errors"

var ErrExecTxError = errors.New("ExecTxError")

type ExecTxError struct {
	message string
}

func (r *ExecTxError) Error() string {
	return r.message
}

func (r *ExecTxError) Unwrap() error {
	return ErrExecTxError
}

// NewExecTxError returns a ExecTxError with the given message.
func NewExecTxError(message string) *ExecTxError {
	return &ExecTxError{
		message: "ExecTxError: " + message,
	}
}
