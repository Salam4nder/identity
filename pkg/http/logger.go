package http

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// ResponseRecorder is a wrapper around http.ResponseWriter that records its
// HTTP status code and body size.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader records the HTTP status code.
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write is a wrapper around http.ResponseWriter.Write that records the
// response body.
func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

// MiddlewareLogger logs HTTP requests. Used for gRPC Gateway.
func MiddlewareLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)

		duration := time.Since(startTime)

		logger := log.Info()

		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http 1.1").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Send()
	})
}
