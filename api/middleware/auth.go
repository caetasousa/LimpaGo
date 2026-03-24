package middleware

import (
	"context"
	"net/http"
	"strings"

	"limpaGo/api/auth"
)

type chaveContexto string

// ChaveUsuarioID é a chave usada para armazenar o ID do usuário no contexto.
const ChaveUsuarioID chaveContexto = "usuario_id"

// AutenticacaoJWT é um middleware que valida o token JWT do header Authorization.
// Extrai o ID do usuário das claims e o armazena no contexto da requisição.
func AutenticacaoJWT(svcToken *auth.ServicoToken) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"codigo":401,"mensagem":"autenticação necessária"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := svcToken.ValidarTokenAcesso(tokenStr)
			if err != nil {
				http.Error(w, `{"codigo":401,"mensagem":"token inválido ou expirado"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ChaveUsuarioID, claims.UsuarioID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ObterUsuarioID recupera o ID do usuário autenticado do contexto.
func ObterUsuarioID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ChaveUsuarioID).(int)
	return id, ok
}
