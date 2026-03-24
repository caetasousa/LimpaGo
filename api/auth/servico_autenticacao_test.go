package auth_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/api/auth"
	erros "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func novoServicoAuth(t *testing.T) *auth.ServicoAutenticacao {
	t.Helper()

	repoUsuario := testutil.NovoRepositorioUsuarioMock()
	repoPerfil := testutil.NovoRepositorioPerfilMock()
	repoCredencial := auth.NovoRepositorioCredencialMock()
	svcUsuario := service.NovoServicoUsuario(repoUsuario, repoPerfil)
	svcToken := auth.NovoServicoToken(auth.ConfiguracaoPadrao())

	return auth.NovoServicoAutenticacao(repoUsuario, repoCredencial, svcUsuario, svcToken)
}

func TestServicoAutenticacao_Registrar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		email       string
		nomeUsuario string
		senha       string
		wantErr     bool
		wantErrIs   error
	}{
		{
			name:        "registro válido",
			email:       "joao@ex.com",
			nomeUsuario: "joao",
			senha:       "Senha123",
			wantErr:     false,
		},
		{
			name:        "senha fraca",
			email:       "joao@ex.com",
			nomeUsuario: "joao",
			senha:       "fraca",
			wantErr:     true,
			wantErrIs:   auth.ErrSenhaFraca,
		},
		{
			name:        "email inválido",
			email:       "",
			nomeUsuario: "joao",
			senha:       "Senha123",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := novoServicoAuth(t)
			ctx := context.Background()

			usuario, tokens, err := svc.Registrar(ctx, tt.email, tt.nomeUsuario, tt.senha)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error; got nil")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("got err %v; want %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if usuario == nil {
				t.Fatal("expected usuario; got nil")
			}
			if tokens == nil {
				t.Fatal("expected tokens; got nil")
			}
			if tokens.TokenAcesso == "" {
				t.Error("expected token_acesso; got empty")
			}
			if tokens.TokenRenovacao == "" {
				t.Error("expected token_renovacao; got empty")
			}
			if tokens.TipoToken != "Bearer" {
				t.Errorf("got tipo_token %q; want %q", tokens.TipoToken, "Bearer")
			}
		})
	}
}

func TestServicoAutenticacao_Registrar_EmailDuplicado(t *testing.T) {
	t.Parallel()

	svc := novoServicoAuth(t)
	ctx := context.Background()

	_, _, err := svc.Registrar(ctx, "joao@ex.com", "joao", "Senha123")
	if err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}

	_, _, err = svc.Registrar(ctx, "joao@ex.com", "joao2", "Senha123")
	if err == nil {
		t.Fatal("expected error on duplicate email; got nil")
	}
	if !errors.Is(err, erros.ErrEmailJaUtilizado) {
		t.Errorf("got err %v; want %v", err, erros.ErrEmailJaUtilizado)
	}
}

func TestServicoAutenticacao_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		email     string
		senha     string
		wantErr   bool
		wantErrIs error
	}{
		{
			name:    "login válido",
			email:   "joao@ex.com",
			senha:   "Senha123",
			wantErr: false,
		},
		{
			name:      "senha errada",
			email:     "joao@ex.com",
			senha:     "SenhaErrada1",
			wantErr:   true,
			wantErrIs: auth.ErrCredenciaisInvalidas,
		},
		{
			name:      "email inexistente",
			email:     "naoexiste@ex.com",
			senha:     "Senha123",
			wantErr:   true,
			wantErrIs: auth.ErrCredenciaisInvalidas,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := novoServicoAuth(t)
			ctx := context.Background()

			// Registrar usuário antes de testar login
			_, _, err := svc.Registrar(ctx, "joao@ex.com", "joao", "Senha123")
			if err != nil {
				t.Fatalf("unexpected error on register: %v", err)
			}

			usuario, tokens, err := svc.Login(ctx, tt.email, tt.senha)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error; got nil")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("got err %v; want %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if usuario == nil {
				t.Fatal("expected usuario; got nil")
			}
			if tokens == nil {
				t.Fatal("expected tokens; got nil")
			}
		})
	}
}

func TestServicoAutenticacao_RenovarToken(t *testing.T) {
	t.Parallel()

	svc := novoServicoAuth(t)
	ctx := context.Background()

	_, tokens, err := svc.Registrar(ctx, "joao@ex.com", "joao", "Senha123")
	if err != nil {
		t.Fatalf("unexpected error on register: %v", err)
	}

	novosTokens, err := svc.RenovarToken(ctx, tokens.TokenRenovacao)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if novosTokens.TokenAcesso == "" {
		t.Error("expected token_acesso; got empty")
	}
	if novosTokens.TipoToken != "Bearer" {
		t.Errorf("got tipo_token %q; want %q", novosTokens.TipoToken, "Bearer")
	}
}

func TestServicoAutenticacao_RenovarToken_Invalido(t *testing.T) {
	t.Parallel()

	svc := novoServicoAuth(t)
	ctx := context.Background()

	_, err := svc.RenovarToken(ctx, "token-invalido")
	if err == nil {
		t.Fatal("expected error; got nil")
	}
	if !errors.Is(err, auth.ErrTokenRenovacaoInvalido) {
		t.Errorf("got err %v; want %v", err, auth.ErrTokenRenovacaoInvalido)
	}
}
