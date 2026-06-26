package middleware

import (
	"net/http"
	"time"

	"github.com/spikycham/finance/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		logger.Infof("%s %s %d %s %s %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start).String(),
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}
