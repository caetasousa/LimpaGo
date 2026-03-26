package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Logger é um middleware que registra método, path, status e duração de cada requisição.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		fmt.Printf("[%s] %s %s %d %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			rw.status,
			time.Since(inicio),
		)
	})
}
