package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	APIKeyHeader = "X-API-KEY"
	ValidAPIKey  = "secret12345"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	message    string
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		message:        "",
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(APIKeyHeader)

		if apiKey == "" || apiKey != ValidAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := newLoggingResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		var message string
		switch wrapped.statusCode {
		case http.StatusOK:
			message = "request processed successfully"
		case http.StatusCreated:
			message = "resource created successfully"
		case http.StatusUnauthorized:
			message = "unauthorized access attempt"
		case http.StatusNotFound:
			message = "resource not found"
		case http.StatusBadRequest:
			message = "bad request"
		case http.StatusMethodNotAllowed:
			message = "method not allowed"
		default:
			message = "request completed"
		}

		log.Printf("%s %s %s %s",
			start.Format("2006-01-02T15:04:05"),
			r.Method,
			r.URL.Path,
			message,
		)
	})
}
