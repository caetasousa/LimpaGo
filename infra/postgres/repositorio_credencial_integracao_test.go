//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/infra/postgres"
)

func TestRepositorioCredencialPG_Salvar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioCredencialPG(db)
	ctx := context.Background()

	usuarioID := inserirUsuario(t, db, "cred@teste.com", "credteste")

	tests := []struct {
		name       string
		usuarioID  int
		senhaHash  string
		wantErr    bool
	}{
		{
			name:      "salvar credencial nova",
			usuarioID: usuarioID,
			senhaHash: "$2a$10$hashficticiodasenha1234",
		},
		{
			name:      "upsert: atualizar credencial existente",
			usuarioID: usuarioID,
			senhaHash: "$2a$10$hashficticionovasenha56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := &auth.Credencial{
				UsuarioID: tt.usuarioID,
				SenhaHash: tt.senhaHash,
			}
			err := repo.Salvar(ctx, cred)
			if tt.wantErr {
				if err == nil {
					t.Fatal("got nil; want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Salvar() unexpected error: %v", err)
			}
			if cred.CriadoEm.IsZero() {
				t.Error("CriadoEm zerado; want preenchido")
			}
		})
	}
}

func TestRepositorioCredencialPG_BuscarPorUsuarioID(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioCredencialPG(db)
	ctx := context.Background()

	usuarioID := inserirUsuario(t, db, "buscarcred@teste.com", "buscarcred")
	cred := &auth.Credencial{UsuarioID: usuarioID, SenhaHash: "$2a$10$hashsalvo"}
	if err := repo.Salvar(ctx, cred); err != nil {
		t.Fatalf("setup: %v", err)
	}

	tests := []struct {
		name      string
		usuarioID int
		wantHash  string
		wantNil   bool
	}{
		{name: "encontrada", usuarioID: usuarioID, wantHash: "$2a$10$hashsalvo"},
		{name: "nao encontrada retorna nil nil", usuarioID: 999999, wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorUsuarioID(ctx, tt.usuarioID)
			if err != nil {
				t.Fatalf("BuscarPorUsuarioID() error: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("got %v; want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil; want credencial")
			}
			if got.SenhaHash != tt.wantHash {
				t.Errorf("SenhaHash = %q; want %q", got.SenhaHash, tt.wantHash)
			}
		})
	}
}
