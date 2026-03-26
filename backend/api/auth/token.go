package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"limpaGo/domain/entity"
)

// ClaimsPersonalizadas contém os dados do usuário embutidos no token JWT.
type ClaimsPersonalizadas struct {
	UsuarioID   int    `json:"usuario_id"`
	Email       string `json:"email"`
	NomeUsuario string `json:"nome_usuario"`
	jwt.RegisteredClaims
}

// ServicoToken gerencia geração e validação de tokens JWT.
type ServicoToken struct {
	config ConfiguracaoJWT
}

// NovoServicoToken cria um novo ServicoToken com a configuração fornecida.
func NovoServicoToken(config ConfiguracaoJWT) *ServicoToken {
	return &ServicoToken{config: config}
}

// GerarTokenAcesso gera um token de acesso JWT de curta duração para o usuário.
func (s *ServicoToken) GerarTokenAcesso(usuario *entity.Usuario) (string, error) {
	agora := time.Now()
	expira := agora.Add(s.config.DuracaoAcesso)

	claims := ClaimsPersonalizadas{
		UsuarioID:   usuario.ID,
		Email:       usuario.Email,
		NomeUsuario: usuario.NomeUsuario,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Emissor,
			Subject:   fmt.Sprintf("%d", usuario.ID),
			IssuedAt:  jwt.NewNumericDate(agora),
			ExpiresAt: jwt.NewNumericDate(expira),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.SegredoAcesso)
}

// GerarTokenRenovacao gera um token de renovação JWT de longa duração para o usuário.
func (s *ServicoToken) GerarTokenRenovacao(usuario *entity.Usuario) (string, error) {
	agora := time.Now()
	expira := agora.Add(s.config.DuracaoRenovacao)

	claims := ClaimsPersonalizadas{
		UsuarioID:   usuario.ID,
		Email:       usuario.Email,
		NomeUsuario: usuario.NomeUsuario,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Emissor,
			Subject:   fmt.Sprintf("%d", usuario.ID),
			IssuedAt:  jwt.NewNumericDate(agora),
			ExpiresAt: jwt.NewNumericDate(expira),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.SegredoRenovacao)
}

// ValidarTokenAcesso valida o token de acesso e retorna as claims se válido.
func (s *ServicoToken) ValidarTokenAcesso(tokenString string) (*ClaimsPersonalizadas, error) {
	return s.validarToken(tokenString, s.config.SegredoAcesso)
}

// ValidarTokenRenovacao valida o token de renovação e retorna as claims se válido.
func (s *ServicoToken) ValidarTokenRenovacao(tokenString string) (*ClaimsPersonalizadas, error) {
	return s.validarToken(tokenString, s.config.SegredoRenovacao)
}

// TempoExpiracaoAcesso retorna o unix timestamp de expiração do próximo token de acesso.
func (s *ServicoToken) TempoExpiracaoAcesso() int64 {
	return time.Now().Add(s.config.DuracaoAcesso).Unix()
}

func (s *ServicoToken) validarToken(tokenString string, segredo []byte) (*ClaimsPersonalizadas, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsPersonalizadas{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return segredo, nil
	})

	if err != nil {
		return nil, ErrTokenInvalido
	}

	claims, ok := token.Claims.(*ClaimsPersonalizadas)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalido
	}

	return claims, nil
}
