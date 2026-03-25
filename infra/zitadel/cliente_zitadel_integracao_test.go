//go:build integration

package zitadel_test

import (
	"context"
	"testing"
)

func TestZitadel_RegistroDeNovoUsuarioViaManagementAPI(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	ctx := context.Background()

	email := emailTeste(t)
	idExterno, err := cliente.RegistrarUsuario(ctx, email, "testuser_reg", "Senha123!")
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}
	if idExterno == "" {
		t.Error("got empty idExterno; want non-empty")
	}
}

func TestZitadel_AutenticacaoComCredenciaisValidasRetornaTokens(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"
	_, err := cliente.RegistrarUsuario(ctx, email, "testuser_auth", senha)
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	tokens, err := cliente.Autenticar(ctx, email, senha)
	if err != nil {
		t.Fatalf("Autenticar() error: %v", err)
	}
	if tokens.TokenAcesso == "" {
		t.Error("got empty token_acesso; want non-empty")
	}
	if tokens.TokenRenovacao == "" {
		t.Error("got empty token_renovacao; want non-empty")
	}
}

func TestZitadel_AutenticacaoComSenhaIncorretaRetornaErro(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	ctx := context.Background()

	email := emailTeste(t)
	_, err := cliente.RegistrarUsuario(ctx, email, "testuser_err", "Senha123!")
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	_, err = cliente.Autenticar(ctx, email, "SenhaErrada999!")
	if err == nil {
		t.Fatal("expected error for wrong password; got nil")
	}
}

func TestZitadel_RenovacaoDeTokenComRefreshTokenValido(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"
	_, err := cliente.RegistrarUsuario(ctx, email, "testuser_renov", senha)
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	tokens, err := cliente.Autenticar(ctx, email, senha)
	if err != nil {
		t.Fatalf("Autenticar() error: %v", err)
	}

	novosTokens, err := cliente.RenovarToken(ctx, tokens.TokenRenovacao)
	if err != nil {
		t.Fatalf("RenovarToken() error: %v", err)
	}
	if novosTokens.TokenAcesso == "" {
		t.Error("got empty token_acesso após renovação; want non-empty")
	}
}

func TestZitadel_RenovacaoDeTokenComRefreshTokenInvalidoRetornaErro(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	ctx := context.Background()

	_, err := cliente.RenovarToken(ctx, "refresh-token-invalido-qualquer")
	if err == nil {
		t.Fatal("expected error for invalid refresh token; got nil")
	}
}

func TestZitadel_ValidacaoDeTokenAcessoViaJWKS(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	svcToken := criarServicoTokenOIDCTeste(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"
	_, err := cliente.RegistrarUsuario(ctx, email, "testuser_jwks", senha)
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	tokens, err := cliente.Autenticar(ctx, email, senha)
	if err != nil {
		t.Fatalf("Autenticar() error: %v", err)
	}

	claims, err := svcToken.ValidarTokenAcesso(tokens.TokenAcesso)
	if err != nil {
		t.Fatalf("ValidarTokenAcesso() error: %v", err)
	}
	if claims.Email != email {
		t.Errorf("got email %q; want %q", claims.Email, email)
	}
}
