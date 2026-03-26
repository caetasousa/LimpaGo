package entity

import (
	"regexp"
	"strings"
	"time"
)

var regexNomeUsuario = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,}$`)

type Usuario struct {
	ID               int
	Email            string
	NomeUsuario      string
	EmailVerificado  bool
	Ativo            bool
	SuperUsuario     bool
	Perfil           *Perfil
	PerfilProfissional  *PerfilProfissional
	PerfilCliente    *PerfilCliente
	CriadoEm        time.Time
	AtualizadoEm    time.Time
}

// NovoUsuario valida e cria um Usuario (sem senha — autenticação é responsabilidade da infraestrutura).
func NovoUsuario(email, nomeUsuario string) (*Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, &ErroValidacao{Campo: "email", Mensagem: "email é obrigatório"}
	}

	if !regexNomeUsuario.MatchString(nomeUsuario) {
		return nil, &ErroValidacao{Campo: "nome_usuario", Mensagem: "deve ter pelo menos 3 caracteres e conter apenas letras, números, underscores ou hífens"}
	}

	return &Usuario{
		Email:       email,
		NomeUsuario: nomeUsuario,
		Ativo:       true,
	}, nil
}

// EProfissional retorna true se o usuário possui um perfil de profissional.
func (u *Usuario) EProfissional() bool {
	return u.PerfilProfissional != nil
}

// ECliente retorna true se o usuário possui um perfil de cliente.
func (u *Usuario) ECliente() bool {
	return u.PerfilCliente != nil
}
