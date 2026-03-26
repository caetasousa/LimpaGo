package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Recuperacao é um middleware que captura panics e retorna 500.
func Recuperacao(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf("[PANIC] %v\n%s\n", rec, debug.Stack())
				http.Error(w, `{"codigo":500,"mensagem":"erro interno do servidor"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
