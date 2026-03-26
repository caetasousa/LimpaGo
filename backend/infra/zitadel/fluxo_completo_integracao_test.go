//go:build integration

package zitadel_test

import (
	"context"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
	"limpaGo/infra/postgres"
)

func criarServicos(t *testing.T) (*auth.ServicoSincronizacao, *auth.ServicoTokenOIDC) {
	t.Helper()
	db := criarBancoTeste(t)
	repoUsuarios := postgres.NovoRepositorioUsuarioPG(db)
	repoPerfis := postgres.NovoRepositorioPerfilPG(db)
	svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)
	sincronizacao := auth.NovoServicoSincronizacao(repoUsuarios, svcUsuario)
	svcToken := criarServicoTokenOIDCTeste(t)
	return sincronizacao, svcToken
}

func TestFluxoCompleto_RegistroNoZitadelSincronizaUsuarioNoBancoLocal(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	sincronizacao, _ := criarServicos(t)
	ctx := context.Background()

	email := emailTeste(t)
	_, err := cliente.RegistrarUsuario(ctx, email, "fluxo_reg", "Senha123!")
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	usuario, err := sincronizacao.SincronizarOuBuscar(ctx, email, "fluxo_reg")
	if err != nil {
		t.Fatalf("SincronizarOuBuscar() error: %v", err)
	}

	if usuario.ID == 0 {
		t.Error("got ID 0; want > 0")
	}
	if usuario.Email != email {
		t.Errorf("got email %q; want %q", usuario.Email, email)
	}
}

func TestFluxoCompleto_LoginRetornaTokensEUsuarioSincronizado(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	sincronizacao, svcToken := criarServicos(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"

	_, err := cliente.RegistrarUsuario(ctx, email, "fluxo_login", senha)
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

	usuario, err := sincronizacao.SincronizarOuBuscar(ctx, claims.Email, claims.NomeUsuario)
	if err != nil {
		t.Fatalf("SincronizarOuBuscar() error: %v", err)
	}
	if usuario.ID == 0 {
		t.Error("got ID 0; want > 0")
	}
}

func TestFluxoCompleto_TokenDoZitadelEValidadoPeloMiddleware(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	sincronizacao, svcToken := criarServicos(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"

	_, err := cliente.RegistrarUsuario(ctx, email, "fluxo_mw", senha)
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}
	_, err = sincronizacao.SincronizarOuBuscar(ctx, email, "fluxo_mw")
	if err != nil {
		t.Fatalf("SincronizarOuBuscar() error: %v", err)
	}

	tokens, err := cliente.Autenticar(ctx, email, senha)
	if err != nil {
		t.Fatalf("Autenticar() error: %v", err)
	}

	claims, err := svcToken.ValidarTokenAcesso(tokens.TokenAcesso)
	if err != nil {
		t.Fatalf("ValidarTokenAcesso() error: %v", err)
	}

	usuarioLocal, err := sincronizacao.SincronizarOuBuscar(ctx, claims.Email, claims.NomeUsuario)
	if err != nil {
		t.Fatalf("SincronizarOuBuscar() error na validação: %v", err)
	}

	// Simula o que o middleware faz: verifica que usuario_id seria injetado no contexto
	ctxComID := context.WithValue(ctx, middleware.ChaveUsuarioID, usuarioLocal.ID)
	idCapturado, ok := middleware.ObterUsuarioID(ctxComID)
	if !ok {
		t.Fatal("ObterUsuarioID retornou false; want true")
	}
	if idCapturado != usuarioLocal.ID {
		t.Errorf("got usuario_id %d; want %d", idCapturado, usuarioLocal.ID)
	}
}

func TestFluxoCompleto_SegundoLoginNaoCriaDuplicataNosBancoLocal(t *testing.T) {
	cliente := criarClienteZitadelTeste(t)
	sincronizacao, _ := criarServicos(t)
	ctx := context.Background()

	email := emailTeste(t)
	senha := "Senha123!"

	_, err := cliente.RegistrarUsuario(ctx, email, "fluxo_dedup", senha)
	if err != nil {
		t.Fatalf("RegistrarUsuario() error: %v", err)
	}

	u1, err := sincronizacao.SincronizarOuBuscar(ctx, email, "fluxo_dedup")
	if err != nil {
		t.Fatalf("primeira sincronização error: %v", err)
	}

	u2, err := sincronizacao.SincronizarOuBuscar(ctx, email, "fluxo_dedup")
	if err != nil {
		t.Fatalf("segunda sincronização error: %v", err)
	}

	if u1.ID != u2.ID {
		t.Errorf("segunda sincronização criou duplicata: got IDs %d e %d", u1.ID, u2.ID)
	}
}
