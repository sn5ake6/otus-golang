package internalhttp

import (
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.statusCode = http.StatusOK

	return w.ResponseWriter.Write(data)
}

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{w, 0}
		start := time.Now()
		next.ServeHTTP(writer, r)
		logger.LogRequest(r, writer.statusCode, time.Since(start))
	})
}
