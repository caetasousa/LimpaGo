//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/infra/postgres"
)

func TestRepositorioUsuarioPG_Salvar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioUsuarioPG(db)
	ctx := context.Background()

	// Inserir um usuario inicial para testar duplicatas
	primeiroUsuario := &entity.Usuario{Email: "joao@exemplo.com", NomeUsuario: "joao123", Ativo: true}
	if err := repo.Salvar(ctx, primeiroUsuario); err != nil {
		t.Fatalf("setup: %v", err)
	}

	tests := []struct {
		name        string
		email       string
		nomeUsuario string
		wantErr     bool
		wantErrIs   error
	}{
		{
			name:        "novo usuario valido",
			email:       "maria@exemplo.com",
			nomeUsuario: "maria456",
		},
		{
			name:        "email duplicado",
			email:       "joao@exemplo.com",
			nomeUsuario: "outronome",
			wantErr:     true,
			wantErrIs:   errosdominio.ErrEmailJaUtilizado,
		},
		{
			name:        "nome_usuario duplicado",
			email:       "carlos@exemplo.com",
			nomeUsuario: "joao123",
			wantErr:     true,
			wantErrIs:   errosdominio.ErrNomeUsuarioJaUtilizado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &entity.Usuario{Email: tt.email, NomeUsuario: tt.nomeUsuario, Ativo: true}
			err := repo.Salvar(ctx, u)

			if tt.wantErr {
				if err == nil {
					t.Fatal("got nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("got %v; want %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("Salvar() unexpected error: %v", err)
			}
			if u.ID == 0 {
				t.Error("ID = 0; want > 0")
			}
			if u.CriadoEm.IsZero() {
				t.Error("CriadoEm zerado; want preenchido")
			}
		})
	}
}

func TestRepositorioUsuarioPG_BuscarPorID(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioUsuarioPG(db)
	ctx := context.Background()

	id := inserirUsuario(t, db, "buscar@id.com", "buscarid")

	tests := []struct {
		name      string
		id        int
		wantEmail string
		wantNil   bool
	}{
		{name: "encontrado", id: id, wantEmail: "buscar@id.com"},
		{name: "nao encontrado retorna nil nil", id: 999999, wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorID(ctx, tt.id)
			if err != nil {
				t.Fatalf("BuscarPorID() error: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("got %v; want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil; want usuario")
			}
			if got.Email != tt.wantEmail {
				t.Errorf("Email = %q; want %q", got.Email, tt.wantEmail)
			}
		})
	}
}

func TestRepositorioUsuarioPG_BuscarPorEmail(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioUsuarioPG(db)
	ctx := context.Background()

	inserirUsuario(t, db, "email@busca.com", "emailbusca")

	tests := []struct {
		name    string
		email   string
		wantNil bool
	}{
		{name: "encontrado", email: "email@busca.com"},
		{name: "nao encontrado", email: "naoexiste@ex.com", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorEmail(ctx, tt.email)
			if err != nil {
				t.Fatalf("BuscarPorEmail() error: %v", err)
			}
			if tt.wantNil && got != nil {
				t.Errorf("got %v; want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Error("got nil; want usuario")
			}
		})
	}
}

func TestRepositorioUsuarioPG_BuscarPorNomeUsuario(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioUsuarioPG(db)
	ctx := context.Background()

	inserirUsuario(t, db, "nome@busca.com", "nomebusca")

	tests := []struct {
		name        string
		nomeUsuario string
		wantNil     bool
	}{
		{name: "encontrado", nomeUsuario: "nomebusca"},
		{name: "nao encontrado", nomeUsuario: "naoexiste", wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorNomeUsuario(ctx, tt.nomeUsuario)
			if err != nil {
				t.Fatalf("BuscarPorNomeUsuario() error: %v", err)
			}
			if tt.wantNil && got != nil {
				t.Errorf("got %v; want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Error("got nil; want usuario")
			}
		})
	}
}
