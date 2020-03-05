package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func httpLoggerMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			lrw := newLoggingResponseWriter(w)

			next.ServeHTTP(lrw, r)

			fields := []zapcore.Field{
				zap.Int("status", lrw.statusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "server.http"),
			}

			zap.L().Info("HTTP Request", fields...)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
