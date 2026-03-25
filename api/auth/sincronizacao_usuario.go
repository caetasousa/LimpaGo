package auth

import (
	"context"
	"errors"
	"strings"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/entity"
	"limpaGo/domain/repository"
	"limpaGo/domain/service"
)

// ServicoSincronizacao provisiona usuários locais a partir de identidades do Zitadel.
// No primeiro acesso, cria o usuário no banco local (just-in-time provisioning).
type ServicoSincronizacao struct {
	usuarios   repository.RepositorioUsuario
	svcUsuario *service.ServicoUsuario
}

// NovoServicoSincronizacao cria um ServicoSincronizacao com as dependências necessárias.
func NovoServicoSincronizacao(
	usuarios repository.RepositorioUsuario,
	svcUsuario *service.ServicoUsuario,
) *ServicoSincronizacao {
	return &ServicoSincronizacao{
		usuarios:   usuarios,
		svcUsuario: svcUsuario,
	}
}

// SincronizarOuBuscar busca o usuário local pelo email.
// Se não existir, cria automaticamente via ServicoUsuario.Registrar.
func (s *ServicoSincronizacao) SincronizarOuBuscar(
	ctx context.Context,
	email string,
	nomeUsuario string,
) (*entity.Usuario, error) {
	usuario, err := s.usuarios.BuscarPorEmail(ctx, email)
	if err == nil && usuario != nil {
		return usuario, nil
	}

	if err != nil && !errors.Is(err, errosdominio.ErrUsuarioNaoEncontrado) {
		return nil, err
	}

	// Usuário não existe — deriva nome de usuário do email se não fornecido
	if nomeUsuario == "" {
		nomeUsuario = derivarNomeUsuario(email)
	}

	return s.svcUsuario.Registrar(ctx, email, nomeUsuario)
}

// derivarNomeUsuario extrai a parte local do email e a usa como nome de usuário base.
func derivarNomeUsuario(email string) string {
	partes := strings.SplitN(email, "@", 2)
	if len(partes) == 0 {
		return "usuario"
	}
	// Remove caracteres inválidos (aceita letras, dígitos, _ e -)
	nome := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, partes[0])
	if len(nome) < 3 {
		nome = nome + "_user"
	}
	return nome
}
