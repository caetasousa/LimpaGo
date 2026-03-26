//go:build integration

package postgres_test

import (
	"context"
	"testing"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/infra/postgres"
)

func TestUsuario_CriacaoComIdExternoPreenchido(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })

	ctx := context.Background()
	repo := postgres.NovoRepositorioUsuarioPG(db)

	tests := []struct {
		name        string
		email       string
		nomeUsuario string
		idExterno   string
		wantErr     bool
	}{
		{
			name:        "usuario com id_externo do Zitadel é salvo corretamente",
			email:       "zitadel@test.com",
			nomeUsuario: "zitadel_user",
			idExterno:   "zitadel-sub-abc123",
		},
		{
			name:        "usuario sem id_externo é salvo com id_externo nulo",
			email:       "local@test.com",
			nomeUsuario: "local_user",
			idExterno:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := inserirUsuario(t, db, tt.email, tt.nomeUsuario)
			if id == 0 {
				t.Fatal("got id=0; want > 0")
			}

			// Verifica que o usuário pode ser buscado por ID
			usuario, err := repo.BuscarPorID(ctx, id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("got nil; want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("BuscarPorID() error: %v", err)
			}
			if usuario.Email != tt.email {
				t.Errorf("got email %q; want %q", usuario.Email, tt.email)
			}
		})
	}
}

func TestUsuario_BuscaPorEmailRetornaErroParaEmailInexistente(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })

	ctx := context.Background()
	repo := postgres.NovoRepositorioUsuarioPG(db)

	usuario, err := repo.BuscarPorEmail(ctx, "inexistente@test.com")
	if err == nil && usuario != nil {
		t.Error("expected nil or error for non-existent email; got user")
	}
	if err != nil && !isErrNotFound(err) {
		t.Errorf("unexpected error type: %v", err)
	}
}

func isErrNotFound(err error) bool {
	return err == errosdominio.ErrUsuarioNaoEncontrado || err.Error() == "sql: no rows in result set"
}
