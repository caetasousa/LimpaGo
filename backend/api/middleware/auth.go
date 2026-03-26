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

// AutenticacaoOIDC é um middleware que valida o token JWT do Zitadel (via JWKS)
// e provisiona o usuário local caso seja o primeiro acesso.
// Injeta o ID interno do usuário no contexto da requisição.
func AutenticacaoOIDC(svcToken *auth.ServicoTokenOIDC, sincronizacao *auth.ServicoSincronizacao) func(http.Handler) http.Handler {
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

			usuario, err := sincronizacao.SincronizarOuBuscar(r.Context(), claims.Email, claims.NomeUsuario)
			if err != nil {
				http.Error(w, `{"codigo":500,"mensagem":"erro ao sincronizar usuário"}`, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ChaveUsuarioID, usuario.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ObterUsuarioID recupera o ID do usuário autenticado do contexto.
func ObterUsuarioID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ChaveUsuarioID).(int)
	return id, ok
}
