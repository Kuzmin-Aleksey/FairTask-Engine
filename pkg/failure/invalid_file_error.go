package failure

import "errors"

type InvalidFileError struct {
	baseError
}

func NewInvalidFileError(msg string) error {
	return InvalidFileError{
		baseError: newBaseError(msg),
	}
}

func (err InvalidFileError) Error() string {
	return "invalid file error: " + err.baseError.Error()
}

func IsInvalidFileError(err error) bool {
	return errors.As(err, new(InvalidFileError))
}
