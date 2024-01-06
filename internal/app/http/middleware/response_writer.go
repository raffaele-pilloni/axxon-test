package middleware

import "net/http"

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newLogResponseWriter(
	responseWriter http.ResponseWriter,
) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: responseWriter}
}

func (l *logResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.size = size

	return size, err
}

func (l *logResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.statusCode = statusCode
}
