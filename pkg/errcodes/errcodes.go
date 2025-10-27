package errcodes

type Code string

func (e Code) String() string {
	return string(e)
}

const (
	ErrUnknown                Code = "unknown error"
	ErrNotFound               Code = "not found"
	ErrContextDone            Code = "context done"
	ErrInvalidRequest         Code = "invalid request"
	ErrDatabaseError          Code = "database error"
	ErrUnauthorized           Code = "unauthorized"
	ErrLocked                 Code = "locked"
	ErrInvalidLoginOrPassword Code = "invalid login or password"
	ErrInvalidFile            Code = "invalid file"
)
