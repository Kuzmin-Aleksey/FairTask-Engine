package server

import (
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/errcodes"
	"FairTask_Engine/pkg/failure"
	"encoding/json"
	"golang.org/x/net/context"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeAndLogErr(ctx context.Context, w http.ResponseWriter, err error) {
	errCode, statusCode := getErrorStatus(err)
	writeJson(ctx, w, ErrorResponse{Error: errCode.String()}, statusCode)

	l := contextx.GetLoggerOrDefault(ctx)
	l.LogAttrs(ctx, slog.LevelError, "error handling request", slog.String("err", err.Error()))
}

func writeJson(ctx context.Context, w http.ResponseWriter, v any, status int) {
	l := contextx.GetLoggerOrDefault(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		l.LogAttrs(ctx, slog.LevelError, "json encode error", slog.String("err", err.Error()))
	}
}

func getErrorStatus(err error) (errcodes.Code, int) {
	switch {
	case failure.IsNotFoundError(err):
		return errcodes.ErrNotFound, http.StatusNotFound
	case failure.IsInvalidRequestError(err):
		return errcodes.ErrInvalidRequest, http.StatusBadRequest
	case failure.IsInvalidFileError(err):
		return errcodes.ErrInvalidFile, http.StatusBadRequest
	case failure.IsUnauthorizedError(err):
		return errcodes.ErrUnauthorized, http.StatusUnauthorized
	case failure.IsLockedError(err):
		return errcodes.ErrLocked, http.StatusLocked
	default:
		return errcodes.ErrUnknown, http.StatusInternalServerError
	}
}
