package middleware

import (
	"log"
	"net/http"
	"time"
)

const (
	KeyHeaderContentType   = "Content-Type"
	FieldHeaderContentType = "application/json"
)

/**
 * Struct.
 */
type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startRequest := time.Now()

			logResponseWriter := newLogResponseWriter(w)

			m.handleNext(next, logResponseWriter, r)

			log.Printf(
				"Request path='%s' method='%s' status='%d' size=%d duration='%s'",
				r.URL.Path,
				r.Method,
				logResponseWriter.statusCode,
				logResponseWriter.size,
				time.Since(startRequest),
			)
		},
	)
}

func (m *Middleware) handleNext(next http.Handler, w http.ResponseWriter, r *http.Request) {
	w.Header().Set(KeyHeaderContentType, FieldHeaderContentType)

	if r.Header.Get(KeyHeaderContentType) != FieldHeaderContentType {
		http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	next.ServeHTTP(w, r)
}
