package failure

import (
	"runtime"
	"strconv"
)

type baseError struct {
	Msg       string
	Initiator string
}

func newBaseError(msg string) baseError {
	pc, _, line, _ := runtime.Caller(2)

	return baseError{
		Msg:       msg,
		Initiator: runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line),
	}
}

func (e baseError) Error() string {
	return e.Initiator + ": " + e.Msg
}
