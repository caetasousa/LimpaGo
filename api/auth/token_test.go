package auth_test

import (
	"errors"
	"testing"
	"time"

	"limpaGo/api/auth"
	"limpaGo/domain/entity"
)

func novoUsuarioTeste(t *testing.T) *entity.Usuario {
	t.Helper()
	u, err := entity.NovoUsuario("teste@ex.com", "testusr")
	if err != nil {
		t.Fatalf("unexpected error creating user: %v", err)
	}
	u.ID = 42
	return u
}

func TestServicoToken_GerarEValidarTokenAcesso(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoPadrao()
	svc := auth.NovoServicoToken(cfg)
	usuario := novoUsuarioTeste(t)

	token, err := svc.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token; got empty string")
	}

	claims, err := svc.ValidarTokenAcesso(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UsuarioID != usuario.ID {
		t.Errorf("got UsuarioID %d; want %d", claims.UsuarioID, usuario.ID)
	}
	if claims.Email != usuario.Email {
		t.Errorf("got Email %q; want %q", claims.Email, usuario.Email)
	}
}

func TestServicoToken_GerarEValidarTokenRenovacao(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoPadrao()
	svc := auth.NovoServicoToken(cfg)
	usuario := novoUsuarioTeste(t)

	token, err := svc.GerarTokenRenovacao(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := svc.ValidarTokenRenovacao(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UsuarioID != usuario.ID {
		t.Errorf("got UsuarioID %d; want %d", claims.UsuarioID, usuario.ID)
	}
}

func TestServicoToken_TokenExpirado(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoPadrao()
	cfg.DuracaoAcesso = -time.Second // já expirado
	svc := auth.NovoServicoToken(cfg)
	usuario := novoUsuarioTeste(t)

	token, err := svc.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.ValidarTokenAcesso(token)
	if err == nil {
		t.Fatal("expected error; got nil")
	}
	if !errors.Is(err, auth.ErrTokenInvalido) {
		t.Errorf("got err %v; want %v", err, auth.ErrTokenInvalido)
	}
}

func TestServicoToken_TokenAdulterado(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoPadrao()
	svc := auth.NovoServicoToken(cfg)
	usuario := novoUsuarioTeste(t)

	token, err := svc.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Adiciona caractere extra para adulterar
	tokenAdulterado := token + "x"

	_, err = svc.ValidarTokenAcesso(tokenAdulterado)
	if err == nil {
		t.Fatal("expected error; got nil")
	}
	if !errors.Is(err, auth.ErrTokenInvalido) {
		t.Errorf("got err %v; want %v", err, auth.ErrTokenInvalido)
	}
}

func TestServicoToken_SegredoErrado(t *testing.T) {
	t.Parallel()

	cfgGerador := auth.ConfiguracaoPadrao()
	cfgGerador.SegredoAcesso = []byte("segredo-a")
	svcGerador := auth.NovoServicoToken(cfgGerador)

	cfgValidador := auth.ConfiguracaoPadrao()
	cfgValidador.SegredoAcesso = []byte("segredo-b")
	svcValidador := auth.NovoServicoToken(cfgValidador)

	usuario := novoUsuarioTeste(t)
	token, err := svcGerador.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svcValidador.ValidarTokenAcesso(token)
	if err == nil {
		t.Fatal("expected error; got nil")
	}
	if !errors.Is(err, auth.ErrTokenInvalido) {
		t.Errorf("got err %v; want %v", err, auth.ErrTokenInvalido)
	}
}

func TestServicoToken_AccessNaoValidaComSegredoRenovacao(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoPadrao()
	svc := auth.NovoServicoToken(cfg)
	usuario := novoUsuarioTeste(t)

	tokenAcesso, err := svc.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Token de acesso não deve ser válido como token de renovação
	_, err = svc.ValidarTokenRenovacao(tokenAcesso)
	if err == nil {
		t.Fatal("expected error; got nil")
	}
}
