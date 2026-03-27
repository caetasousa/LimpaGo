package auth

import "errors"

var (
	// ErrCredenciaisInvalidas é retornado quando email ou senha estão incorretos.
	// Mensagem genérica para não revelar se o email existe (prevenção de enumeração).
	ErrCredenciaisInvalidas = errors.New("email ou senha incorretos")

	// ErrSenhaFraca é retornado quando a senha não atende os requisitos mínimos.
	ErrSenhaFraca = errors.New("senha deve ter no mínimo 8 caracteres, incluindo letra maiúscula e número")

	// ErrTokenInvalido é retornado quando o token de acesso é inválido ou expirado.
	ErrTokenInvalido = errors.New("token inválido ou expirado")

	// ErrTokenRenovacaoInvalido é retornado quando o token de renovação é inválido ou expirado.
	ErrTokenRenovacaoInvalido = errors.New("token de renovação inválido ou expirado")

	// ErrUsuarioInativo é retornado quando a conta do usuário está desativada.
	ErrUsuarioInativo = errors.New("conta de usuário desativada")

)
