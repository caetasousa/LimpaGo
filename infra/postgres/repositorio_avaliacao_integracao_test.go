//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
	"limpaGo/infra/postgres"
)

func TestRepositorioAvaliacaoPG_Salvar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax@aval.com", "faxaval")
	clienteID := inserirUsuario(t, db, "cli@aval.com", "cliaval")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico Aval")

	tests := []struct {
		name        string
		clienteID   int
		limpezaID   int
		nota        valueobject.Nota
		wantErr     bool
		wantErrIs   error
	}{
		{
			name:      "salvar avaliacao valida",
			clienteID: clienteID,
			limpezaID: limpezaID,
			nota:      4,
		},
		{
			name:      "avaliacao duplicada retorna ErrAvaliacaoDuplicada",
			clienteID: clienteID,
			limpezaID: limpezaID,
			nota:      5,
			wantErr:   true,
			wantErrIs: errosdominio.ErrAvaliacaoDuplicada,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := entity.NovaAvaliacao(tt.limpezaID, faxineiroID, tt.clienteID, tt.nota, "Ótimo serviço")
			err := repo.Salvar(ctx, a)

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
			if a.ID == 0 {
				t.Error("ID = 0; want > 0")
			}
		})
	}
}

func TestRepositorioAvaliacaoPG_BuscarPorClienteELimpeza(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax2@aval.com", "fax2aval")
	clienteID := inserirUsuario(t, db, "cli2@aval.com", "cli2aval")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico2 Aval")

	a := entity.NovaAvaliacao(limpezaID, faxineiroID, clienteID, 5, "Perfeito")
	if err := repo.Salvar(ctx, a); err != nil {
		t.Fatalf("setup: %v", err)
	}

	tests := []struct {
		name      string
		clienteID int
		limpezaID int
		wantErr   error
		wantNil   bool
	}{
		{name: "encontrada", clienteID: clienteID, limpezaID: limpezaID},
		{
			name:      "nao encontrada",
			clienteID: 999999,
			limpezaID: limpezaID,
			wantErr:   errosdominio.ErrAvaliacaoNaoEncontrada,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorClienteELimpeza(ctx, tt.clienteID, tt.limpezaID)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v; want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if got == nil {
				t.Fatal("got nil; want avaliacao")
			}
			if got.Nota != 5 {
				t.Errorf("Nota = %d; want 5", got.Nota)
			}
		})
	}
}

func TestRepositorioAvaliacaoPG_ListarPorFaxineiro(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax3@aval.com", "fax3aval")
	cli1 := inserirUsuario(t, db, "cli3a@aval.com", "cli3a")
	cli2 := inserirUsuario(t, db, "cli3b@aval.com", "cli3b")
	limp1 := inserirLimpeza(t, db, faxineiroID, "Aval1")
	limp2 := inserirLimpeza(t, db, faxineiroID, "Aval2")

	repo.Salvar(ctx, entity.NovaAvaliacao(limp1, faxineiroID, cli1, 4, "Bom"))
	repo.Salvar(ctx, entity.NovaAvaliacao(limp2, faxineiroID, cli2, 5, "Excelente"))

	got, err := repo.ListarPorFaxineiro(ctx, faxineiroID)
	if err != nil {
		t.Fatalf("ListarPorFaxineiro() error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d; want 2", len(got))
	}
}

func TestRepositorioAvaliacaoPG_BuscarAgregadoPorFaxineiro(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax4@aval.com", "fax4aval")
	cli1 := inserirUsuario(t, db, "cli4a@aval.com", "cli4a")
	cli2 := inserirUsuario(t, db, "cli4b@aval.com", "cli4b")
	limp1 := inserirLimpeza(t, db, faxineiroID, "Aval3")
	limp2 := inserirLimpeza(t, db, faxineiroID, "Aval4")

	repo.Salvar(ctx, entity.NovaAvaliacao(limp1, faxineiroID, cli1, 4, ""))
	repo.Salvar(ctx, entity.NovaAvaliacao(limp2, faxineiroID, cli2, 2, ""))

	t.Run("agregado com avaliacoes calcula media", func(t *testing.T) {
		ag, err := repo.BuscarAgregadoPorFaxineiro(ctx, faxineiroID)
		if err != nil {
			t.Fatalf("BuscarAgregadoPorFaxineiro() error: %v", err)
		}
		if ag.TotalAvaliacoes != 2 {
			t.Errorf("TotalAvaliacoes = %d; want 2", ag.TotalAvaliacoes)
		}
		if ag.MediaNota != 3.0 {
			t.Errorf("MediaNota = %v; want 3.0", ag.MediaNota)
		}
	})

	t.Run("faxineiro sem avaliacoes retorna zerado", func(t *testing.T) {
		semAvaliacaoID := inserirUsuario(t, db, "semav@aval.com", "semavaliacao")
		ag, err := repo.BuscarAgregadoPorFaxineiro(ctx, semAvaliacaoID)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if ag.TotalAvaliacoes != 0 {
			t.Errorf("TotalAvaliacoes = %d; want 0", ag.TotalAvaliacoes)
		}
		if ag.MediaNota != 0.0 {
			t.Errorf("MediaNota = %v; want 0.0", ag.MediaNota)
		}
	})
}
