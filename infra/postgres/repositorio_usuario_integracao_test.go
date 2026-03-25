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

func TestCadastro_RegistrarNovoUsuarioNaPlataforma(t *testing.T) {
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
			name:        "usuario com email e nome unicos e registrado com sucesso",
			email:       "maria@exemplo.com",
			nomeUsuario: "maria456",
		},
		{
			name:        "sistema rejeita registro com email ja cadastrado",
			email:       "joao@exemplo.com",
			nomeUsuario: "outronome",
			wantErr:     true,
			wantErrIs:   errosdominio.ErrEmailJaUtilizado,
		},
		{
			name:        "sistema rejeita registro com nome de usuario ja existente",
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

func TestCadastro_BuscarUsuarioPorIDRetornaDadosCorretos(t *testing.T) {
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
		{name: "usuario existente e encontrado pelo ID", id: id, wantEmail: "buscar@id.com"},
		{name: "ID inexistente retorna nulo sem erro", id: 999999, wantNil: true},
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

func TestCadastro_BuscarUsuarioPorEmailParaLogin(t *testing.T) {
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
		{name: "email cadastrado retorna usuario correspondente", email: "email@busca.com"},
		{name: "email nao cadastrado retorna nulo", email: "naoexiste@ex.com", wantNil: true},
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

func TestCadastro_BuscarUsuarioPorNomeDeUsuario(t *testing.T) {
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
		{name: "nome de usuario existente retorna usuario", nomeUsuario: "nomebusca"},
		{name: "nome de usuario inexistente retorna nulo", nomeUsuario: "naoexiste", wantNil: true},
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
