package middlewarex

import (
	"FairTask_Engine/pkg/logx"
	"bytes"
	"cmp"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji/web/mutil"
)

var acceptContentTypes = map[string]bool{
	"application/json": true,
	"application/xml":  true,
	"text/xml":         true,
}

// The trouble with optional interfaces:
// https://blog.merovius.de/posts/2017-07-30-the-trouble-with-optional-interfaces/
// https://medium.com/@cep21/interface-wrapping-method-erasure-c523b3549912
func ResponseLogging(
	sensitiveDataMasker logx.SensitiveDataMaskerInterface,
	logFieldMaxLen int,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.RequestURI, "image") ||
				strings.Contains(r.RequestURI, "static") ||
				strings.Contains(r.RequestURI, "ico") {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			start := time.Now()
			lw := mutil.WrapWriter(w)

			var buf bytes.Buffer

			lw.Tee(&buf)

			next.ServeHTTP(lw, r)

			responseHeaders, err := responseHeaders(w)
			if err != nil {
				logger(ctx).Error("responseHeaders", logx.Error(err))
			}

			dump := buf.Bytes()

			if len(dump) > logFieldMaxLen {
				dump = dump[:logFieldMaxLen]
			}

			// Если в хэндлере принудительно не установлен статус, то
			// lw.Status() будет возвращать 0 (упоминание этого есть в
			// документации). Поэтому устанавливаем статус 200 вручную.
			status := cmp.Or(lw.Status(), http.StatusOK)

			logArgs := []any{
				slog.Int(logx.FieldResponseStatus, status),
				slog.Int64(logx.FieldDurationMs, time.Since(start).Milliseconds()),
				slog.String(logx.FieldResponseHeaders, string(sensitiveDataMasker.Mask(responseHeaders))),
			}

			if acceptContentTypes[w.Header().Get("Content-Type")] {
				logArgs = append(logArgs, slog.String(logx.FieldResponseBody, string(sensitiveDataMasker.Mask(dump))))
			}

			logger(ctx).Info(logx.FieldHTTPResponse, logArgs...)
		})
	}
}

func responseHeaders(w http.ResponseWriter) ([]byte, error) {
	var buf bytes.Buffer

	if err := w.Header().WriteSubset(&buf, nil); err != nil {
		return nil, fmt.Errorf("header.WriteSubset: %w", err)
	}

	return buf.Bytes(), nil
}
