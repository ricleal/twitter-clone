package entities

type InvalidEmailError struct {
	Email string
}

func NewInvalidEmailError(email string) *InvalidEmailError {
	return &InvalidEmailError{email}
}

func (e *InvalidEmailError) Error() string {
	return "invalid email address"
}
