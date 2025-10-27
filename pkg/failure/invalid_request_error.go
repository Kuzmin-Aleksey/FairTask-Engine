package failure

import "errors"

type InvalidRequestError struct {
	baseError
}

func NewInvalidRequestError(msg string) error {
	return InvalidRequestError{
		baseError: newBaseError(msg),
	}
}

func (err InvalidRequestError) Error() string {
	return "invalid request error: " + err.baseError.Error()
}

func IsInvalidRequestError(err error) bool {
	return errors.As(err, new(InvalidRequestError))
}
