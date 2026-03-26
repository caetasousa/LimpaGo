package auth

import (
	"context"
	"errors"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/entity"
)

// ServicoAutenticacao gerencia registro, login e renovação de tokens via Zitadel.
type ServicoAutenticacao struct {
	clienteZitadel *ClienteZitadel
	sincronizacao  *ServicoSincronizacao
	svcToken       *ServicoTokenOIDC
}

// NovoServicoAutenticacao cria um ServicoAutenticacao que delega ao Zitadel.
func NovoServicoAutenticacao(
	clienteZitadel *ClienteZitadel,
	sincronizacao *ServicoSincronizacao,
	svcToken *ServicoTokenOIDC,
) *ServicoAutenticacao {
	return &ServicoAutenticacao{
		clienteZitadel: clienteZitadel,
		sincronizacao:  sincronizacao,
		svcToken:       svcToken,
	}
}

// Registrar cria um novo usuário no Zitadel e sincroniza no banco local.
// Retorna o usuário local e o par de tokens emitido pelo Zitadel.
func (s *ServicoAutenticacao) Registrar(ctx context.Context, email, nomeUsuario, senha string) (*entity.Usuario, *ParTokens, error) {
	// 1. Registra o usuário no Zitadel (Management API)
	_, err := s.clienteZitadel.RegistrarUsuario(ctx, email, nomeUsuario, senha)
	if err != nil {
		if errors.Is(err, ErrEmailJaCadastradoNoIdP) {
			return nil, nil, errosdominio.ErrEmailJaUtilizado
		}
		return nil, nil, err
	}

	// 2. Obtém tokens via login (Resource Owner Password Grant)
	tokens, err := s.clienteZitadel.Autenticar(ctx, email, senha)
	if err != nil {
		return nil, nil, err
	}

	// 3. Sincroniza o usuário no banco local
	usuario, err := s.sincronizacao.SincronizarOuBuscar(ctx, email, nomeUsuario)
	if err != nil {
		return nil, nil, err
	}

	return usuario, zitadelParaParTokens(tokens), nil
}

// Login autentica o usuário via Zitadel e sincroniza no banco local.
func (s *ServicoAutenticacao) Login(ctx context.Context, email, senha string) (*entity.Usuario, *ParTokens, error) {
	// 1. Autentica via Zitadel
	tokens, err := s.clienteZitadel.Autenticar(ctx, email, senha)
	if err != nil {
		if errors.Is(err, ErrCredenciaisInvalidas) {
			return nil, nil, ErrCredenciaisInvalidas
		}
		return nil, nil, err
	}

	// 2. Extrai email das claims do token para sincronizar
	claims, err := s.svcToken.ValidarTokenAcesso(tokens.TokenAcesso)
	emailParaSinc := email
	nomeParaSinc := ""
	if err == nil && claims != nil {
		if claims.Email != "" {
			emailParaSinc = claims.Email
		}
		nomeParaSinc = claims.NomeUsuario
	}

	// 3. Sincroniza o usuário no banco local
	usuario, err := s.sincronizacao.SincronizarOuBuscar(ctx, emailParaSinc, nomeParaSinc)
	if err != nil {
		return nil, nil, err
	}

	if !usuario.Ativo {
		return nil, nil, ErrUsuarioInativo
	}

	return usuario, zitadelParaParTokens(tokens), nil
}

// RenovarToken usa um refresh_token do Zitadel para obter novos tokens.
func (s *ServicoAutenticacao) RenovarToken(ctx context.Context, tokenRenovacao string) (*ParTokens, error) {
	tokens, err := s.clienteZitadel.RenovarToken(ctx, tokenRenovacao)
	if err != nil {
		if errors.Is(err, ErrCredenciaisInvalidas) {
			return nil, ErrTokenRenovacaoInvalido
		}
		return nil, ErrTokenRenovacaoInvalido
	}

	return zitadelParaParTokens(tokens), nil
}

func zitadelParaParTokens(t *TokensZitadel) *ParTokens {
	return &ParTokens{
		TokenAcesso:    t.TokenAcesso,
		TokenRenovacao: t.TokenRenovacao,
		TipoToken:      "Bearer",
		ExpiraEm:       t.ExpiraEm,
	}
}
