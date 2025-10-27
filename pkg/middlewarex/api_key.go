package middlewarex

import (
	"FairTask_Engine/pkg/failure"
	"FairTask_Engine/pkg/logx"
	"net/http"
	"strings"
)

func WithApiKey(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			typeAndKey := strings.Split(r.Header.Get("Authorization"), " ")
			if len(typeAndKey) != 2 {
				logger(ctx).Warn("api key", logx.Error(failure.NewUnauthorizedError("invalid authorization header")))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if typeAndKey[0] != "ApiKey" {
				logger(ctx).Warn("api key", logx.Error(failure.NewUnauthorizedError("invalid auth type")))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			key := typeAndKey[1]
			if key != apiKey {
				logger(ctx).Warn("api key", logx.Error(failure.NewUnauthorizedError("invalid api key")))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
