package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"taskmanager/internal/config"
	"time"
)

// Embedded ResponseWriter to get status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// Redefined WriteHeader method that write status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func APIKeyMiddleware(next http.Handler) http.Handler {
	conf := *config.GetAuthMiddlewareConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(conf.ApiKeyHeader)

		if apiKey == "" || apiKey != conf.ValidAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newLoggingResponseWriter(w)

		next.ServeHTTP(rw, r)
		log.Printf("%s at %s with status code: %d, Server Time: %s", r.Method, r.RequestURI, rw.statusCode, start.Format(time.RFC3339))
	})
}
