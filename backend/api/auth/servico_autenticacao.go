package auth

import (
	"context"
	"errors"
	"unicode"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/entity"
	"limpaGo/domain/repository"
	"limpaGo/domain/service"

	"golang.org/x/crypto/bcrypt"
)

// ServicoAutenticacao gerencia registro, login e renovação de tokens via bcrypt + JWT local.
type ServicoAutenticacao struct {
	usuarios    repository.RepositorioUsuario
	credenciais RepositorioCredencial
	svcUsuario  *service.ServicoUsuario
	svcToken    *ServicoToken
}

// NovoServicoAutenticacao cria um ServicoAutenticacao com autenticação local.
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

// Registrar valida a senha, cria o usuário, salva o hash bcrypt e retorna tokens HMAC.
func (s *ServicoAutenticacao) Registrar(ctx context.Context, email, nomeUsuario, senha string) (*entity.Usuario, *ParTokens, error) {
	if err := validarForcaSenha(senha); err != nil {
		return nil, nil, err
	}

	usuario, err := s.svcUsuario.Registrar(ctx, email, nomeUsuario)
	if err != nil {
		if errors.Is(err, errosdominio.ErrEmailJaUtilizado) {
			return nil, nil, errosdominio.ErrEmailJaUtilizado
		}
		return nil, nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	cred := &Credencial{UsuarioID: usuario.ID, SenhaHash: string(hash)}
	if err := s.credenciais.Salvar(ctx, cred); err != nil {
		return nil, nil, err
	}

	tokens, err := s.gerarTokens(usuario)
	if err != nil {
		return nil, nil, err
	}

	return usuario, tokens, nil
}

// Login verifica email, compara o hash bcrypt e retorna tokens HMAC.
func (s *ServicoAutenticacao) Login(ctx context.Context, email, senha string) (*entity.Usuario, *ParTokens, error) {
	usuario, err := s.usuarios.BuscarPorEmail(ctx, email)
	if err != nil || usuario == nil {
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

	tokens, err := s.gerarTokens(usuario)
	if err != nil {
		return nil, nil, err
	}

	return usuario, tokens, nil
}

// RenovarToken valida o token de renovação HMAC e gera um novo par de tokens.
func (s *ServicoAutenticacao) RenovarToken(ctx context.Context, tokenRenovacao string) (*ParTokens, error) {
	claims, err := s.svcToken.ValidarTokenRenovacao(tokenRenovacao)
	if err != nil {
		return nil, ErrTokenRenovacaoInvalido
	}

	usuario, err := s.usuarios.BuscarPorID(ctx, claims.UsuarioID)
	if err != nil {
		return nil, ErrTokenRenovacaoInvalido
	}

	return s.gerarTokens(usuario)
}

func (s *ServicoAutenticacao) gerarTokens(usuario *entity.Usuario) (*ParTokens, error) {
	tokenAcesso, err := s.svcToken.GerarTokenAcesso(usuario)
	if err != nil {
		return nil, err
	}

	tokenRenovacao, err := s.svcToken.GerarTokenRenovacao(usuario)
	if err != nil {
		return nil, err
	}

	return &ParTokens{
		TokenAcesso:    tokenAcesso,
		TokenRenovacao: tokenRenovacao,
		TipoToken:      "Bearer",
		ExpiraEm:       s.svcToken.TempoExpiracaoAcesso(),
	}, nil
}

// validarForcaSenha exige mínimo 8 caracteres, ao menos uma letra maiúscula e um dígito.
func validarForcaSenha(senha string) error {
	if len(senha) < 8 {
		return ErrSenhaFraca
	}
	temMaiuscula := false
	temDigito := false
	for _, r := range senha {
		if unicode.IsUpper(r) {
			temMaiuscula = true
		}
		if unicode.IsDigit(r) {
			temDigito = true
		}
	}
	if !temMaiuscula || !temDigito {
		return ErrSenhaFraca
	}
	return nil
}
