package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// ClaimsZitadel contém as claims JWT emitidas pelo Zitadel via OIDC.
type ClaimsZitadel struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerificado   bool   `json:"email_verified"`
	NomeUsuario       string `json:"preferred_username"`
	jwt.RegisteredClaims
}

// ServicoTokenOIDC valida tokens JWT emitidos pelo Zitadel usando JWKS (RS256).
type ServicoTokenOIDC struct {
	emissor  string
	clientID string
	keyfunc  keyfunc.Keyfunc
}

// NovoServicoTokenOIDC cria um ServicoTokenOIDC que busca e cacheia as chaves JWKS do Zitadel.
func NovoServicoTokenOIDC(cfg ConfiguracaoZitadel) (*ServicoTokenOIDC, error) {
	jwksURL := fmt.Sprintf("%s/oauth/v2/keys", cfg.URL)

	kf, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar JWKS: %w", err)
	}

	return &ServicoTokenOIDC{
		emissor:  cfg.Emissor,
		clientID: cfg.ClientID,
		keyfunc:  kf,
	}, nil
}

// NovoServicoTokenOIDCMock cria um ServicoTokenOIDC de desenvolvimento sem validação JWKS real.
// Usado quando ZITADEL_URL não está configurada (modo in-memory).
func NovoServicoTokenOIDCMock() *ServicoTokenOIDC {
	return &ServicoTokenOIDC{
		emissor:  "mock",
		clientID: "mock",
		keyfunc:  nil,
	}
}

// ValidarTokenAcesso valida um access_token JWT emitido pelo Zitadel e retorna as claims.
func (s *ServicoTokenOIDC) ValidarTokenAcesso(tokenString string) (*ClaimsZitadel, error) {
	if s.keyfunc == nil {
		return nil, ErrTokenInvalido
	}

	claims := &ClaimsZitadel{}
	token, err := jwt.ParseWithClaims(tokenString, claims, s.keyfunc.Keyfunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenInvalido
		}
		return nil, ErrTokenInvalido
	}

	if !token.Valid {
		return nil, ErrTokenInvalido
	}

	if s.emissor != "mock" {
		issuer, err := claims.GetIssuer()
		if err != nil || issuer != s.emissor {
			return nil, ErrTokenInvalido
		}
	}

	return claims, nil
}

// TempoExpiracaoAcesso retorna o unix timestamp de expiração do próximo token de acesso.
// Usa 15 minutos como duração padrão (alinhado com o Zitadel).
func (s *ServicoTokenOIDC) TempoExpiracaoAcesso() int64 {
	return time.Now().Add(15 * time.Minute).Unix()
}
