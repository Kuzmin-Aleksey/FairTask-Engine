package middlewarex

import (
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"log/slog"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID := contextx.GetTraceId(ctx)

		ctx = contextx.WithLogger(
			ctx,
			logger(ctx).With(
				logx.Stringer(logx.FieldTraceID, traceID),
				logx.Stringer(logx.FieldURL, r.URL),
				slog.String(logx.FieldHTTPMethod, r.Method),
				slog.String(logx.FieldIP, r.RemoteAddr),
			),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
