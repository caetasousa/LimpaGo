// Package middleware contém os middlewares HTTP da API Phresh.
package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type chaveContexto string

// ChaveUsuarioID é a chave usada para armazenar o ID do usuário no contexto.
const ChaveUsuarioID chaveContexto = "usuario_id"

// Autenticacao é um middleware placeholder que extrai o ID do usuário do header X-User-ID.
// TODO: substituir por validação JWT para produção.
func Autenticacao(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerID := r.Header.Get("X-User-ID")
		if headerID == "" {
			http.Error(w, `{"codigo":401,"mensagem":"autenticação necessária"}`, http.StatusUnauthorized)
			return
		}

		usuarioID, err := strconv.Atoi(headerID)
		if err != nil || usuarioID <= 0 {
			http.Error(w, `{"codigo":401,"mensagem":"X-User-ID inválido"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ChaveUsuarioID, usuarioID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ObterUsuarioID recupera o ID do usuário autenticado do contexto.
func ObterUsuarioID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ChaveUsuarioID).(int)
	return id, ok
}
