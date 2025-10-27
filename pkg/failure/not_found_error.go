package failure

import (
	"errors"
)

type NotFoundError struct {
	baseError
}

func NewNotFoundError(msg string) error {
	return NotFoundError{
		baseError: newBaseError(msg),
	}
}

func (err NotFoundError) Error() string {
	return "not found: " + err.baseError.Error()
}

func IsNotFoundError(err error) bool {
	return errors.As(err, new(NotFoundError))
}
