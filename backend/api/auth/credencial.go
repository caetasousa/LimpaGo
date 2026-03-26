package auth

import "time"

// Credencial armazena o hash da senha de um usuário.
// Intencionalmente separada do domínio — autenticação é responsabilidade da infraestrutura.
type Credencial struct {
	UsuarioID    int
	SenhaHash    string
	CriadoEm    time.Time
	AtualizadoEm time.Time
}
