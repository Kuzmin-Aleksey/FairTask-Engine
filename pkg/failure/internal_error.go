package failure

import "errors"

type InternalError struct {
	baseError
}

func NewInternalError(msg string) error {
	return InternalError{
		baseError: newBaseError(msg),
	}
}

func (err InternalError) Error() string {
	return "internal error: " + err.baseError.Error()
}

func IsInternalError(err error) bool {
	return errors.As(err, new(InternalError))
}
