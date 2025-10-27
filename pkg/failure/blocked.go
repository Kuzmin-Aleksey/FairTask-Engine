package failure

import "errors"

type LockedError struct {
	baseError
}

func NewLockedError(msg string) error {
	return LockedError{
		baseError: newBaseError(msg),
	}
}

func (err LockedError) Error() string {
	return "locked: " + err.baseError.Error()
}

func IsLockedError(err error) bool {
	return errors.As(err, new(LockedError))
}
