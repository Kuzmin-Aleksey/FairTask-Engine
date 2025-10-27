package logx

import (
	"fmt"
	"github.com/lmittmann/tint"
	"log/slog"
)

var Error = tint.Err //nolint:gochecknoglobals

func Stringer(name string, value fmt.Stringer) slog.Attr {
	return slog.String(name, value.String())
}
