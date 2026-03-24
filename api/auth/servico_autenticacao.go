package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"limpaGo/domain/entity"
	"limpaGo/domain/repository"
	"limpaGo/domain/service"
)

const custoHash = 12

// ServicoAutenticacao gerencia registro, login e renovação de tokens.
type ServicoAutenticacao struct {
	usuarios    repository.RepositorioUsuario
	credenciais RepositorioCredencial
	svcUsuario  *service.ServicoUsuario
	svcToken    *ServicoToken
}

// NovoServicoAutenticacao cria um novo ServicoAutenticacao com as dependências necessárias.
func NovoServicoAutenticacao(
	usuarios repository.RepositorioUsuario,
	credenciais RepositorioCredencial,
	svcUsuario *service.ServicoUsuario,
	svcToken *ServicoToken,
) *ServicoAutenticacao {
	return &ServicoAutenticacao{
		usuarios:    usuarios,
		credenciais: credenciais,
		svcUsuario:  svcUsuario,
		svcToken:    svcToken,
	}
}

// Registrar cria um novo usuário com senha e retorna um par de tokens JWT.
func (s *ServicoAutenticacao) Registrar(ctx context.Context, email, nomeUsuario, senha string) (*entity.Usuario, *ParTokens, error) {
	if err := ValidarForcaSenha(senha); err != nil {
		return nil, nil, err
	}

	usuario, err := s.svcUsuario.Registrar(ctx, email, nomeUsuario)
	if err != nil {
		return nil, nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), custoHash)
	if err != nil {
		return nil, nil, err
	}

	cred := &Credencial{
		UsuarioID: usuario.ID,
		SenhaHash: string(hash),
	}
	if err := s.credenciais.Salvar(ctx, cred); err != nil {
		return nil, nil, err
	}

	tokens, err := s.gerarParTokens(usuario)
	if err != nil {
		return nil, nil, err
	}

	return usuario, tokens, nil
}

// Login autentica um usuário por email e senha e retorna um par de tokens JWT.
func (s *ServicoAutenticacao) Login(ctx context.Context, email, senha string) (*entity.Usuario, *ParTokens, error) {
	usuario, err := s.usuarios.BuscarPorEmail(ctx, email)
	if err != nil || usuario == nil {
		// Mensagem genérica para prevenir enumeração de emails
		return nil, nil, ErrCredenciaisInvalidas
	}

	if !usuario.Ativo {
		return nil, nil, ErrUsuarioInativo
	}

	cred, err := s.credenciais.BuscarPorUsuarioID(ctx, usuario.ID)
	if err != nil {
		return nil, nil, ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.SenhaHash), []byte(senha)); err != nil {
		return nil, nil, ErrCredenciaisInvalidas
	}

	tokens, err := s.gerarParTokens(usuario)
	if err != nil {
		return nil, nil, err
	}

	return usuario, tokens, nil
}

// RenovarToken valida um token de renovação e retorna um novo par de tokens JWT.
func (s *ServicoAutenticacao) RenovarToken(ctx context.Context, tokenRenovacao string) (*ParTokens, error) {
	claims, err := s.svcToken.ValidarTokenRenovacao(tokenRenovacao)
	if err != nil {
		return nil, ErrTokenRenovacaoInvalido
	}

	usuario, err := s.usuarios.BuscarPorID(ctx, claims.UsuarioID)
	if err != nil {
		return nil, ErrTokenRenovacaoInvalido
	}

	if !usuario.Ativo {
		return nil, ErrUsuarioInativo
	}

	return s.gerarParTokens(usuario)
}

func (s *ServicoAutenticacao) gerarParTokens(usuario *entity.Usuario) (*ParTokens, error) {
	tokenAcesso, err := s.svcToken.GerarTokenAcesso(usuario)
	if err != nil {
		return nil, err
	}

	tokenRenovacao, err := s.svcToken.GerarTokenRenovacao(usuario)
	if err != nil {
		return nil, err
	}

	expiraEm := s.svcToken.TempoExpiracaoAcesso()

	return &ParTokens{
		TokenAcesso:    tokenAcesso,
		TokenRenovacao: tokenRenovacao,
		TipoToken:      "Bearer",
		ExpiraEm:       expiraEm,
	}, nil
}
