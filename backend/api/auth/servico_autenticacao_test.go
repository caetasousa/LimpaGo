package auth_test

import (
	"context"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func novoServicoAuth(t *testing.T) *auth.ServicoAutenticacao {
	t.Helper()
	repoUsuarios := testutil.NovoRepositorioUsuarioMock()
	repoPerfis := testutil.NovoRepositorioPerfilMock()
	svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)
	repoCredenciais := auth.NovoRepositorioCredencialMock()
	cfgJWT := auth.ConfiguracaoPadrao()
	svcToken := auth.NovoServicoToken(cfgJWT)
	return auth.NovoServicoAutenticacao(repoUsuarios, repoCredenciais, svcUsuario, svcToken)
}

func TestServicoAutenticacao_NovoServicoAutenticacao(t *testing.T) {
	t.Parallel()
	svc := novoServicoAuth(t)
	if svc == nil {
		t.Fatal("expected non-nil ServicoAutenticacao; got nil")
	}
}

func TestServicoAutenticacao_RenovarTokenComTokenInvalido(t *testing.T) {
	t.Parallel()
	svc := novoServicoAuth(t)

	_, err := svc.RenovarToken(context.Background(), "refresh-token-invalido")
	if err == nil {
		t.Error("expected error for invalid refresh token; got nil")
	}
}
