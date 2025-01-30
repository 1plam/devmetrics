package logger

import (
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Length     int64
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.Length += int64(n)
	return n, err
}

func HTTPMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create custom response writer to capture status code
			rw := &ResponseWriter{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			}

			// Create request-specific logger with request ID
			requestLogger := logger.With(
				String("request_id", r.Header.Get("X-Request-ID")),
				String("method", r.Method),
				String("path", r.URL.Path),
				String("remote_addr", r.RemoteAddr),
				String("user_agent", r.UserAgent()),
			)

			// Log request
			requestLogger.Info("Incoming request")

			// Call next handler
			next.ServeHTTP(rw, r)

			// Log response
			requestLogger.Info("Request completed",
				Int("status", rw.StatusCode),
				Int64("bytes", rw.Length),
				Duration("duration", time.Since(start)),
			)
		})
	}
}
