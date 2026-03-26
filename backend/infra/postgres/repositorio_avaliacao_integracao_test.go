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

func TestAvaliacao_ClienteAvaliaServicoAposConclusao(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax@aval.com", "faxaval")
	clienteID := inserirUsuario(t, db, "cli@aval.com", "cliaval")
	limpezaID := inserirLimpeza(t, db, profissionalID, "Servico Aval")

	tests := []struct {
		name      string
		clienteID int
		limpezaID int
		nota      valueobject.Nota
		wantErr   bool
		wantErrIs error
	}{
		{
			name:      "cliente registra avaliacao com nota e comentario",
			clienteID: clienteID,
			limpezaID: limpezaID,
			nota:      4,
		},
		{
			name:      "sistema rejeita avaliacao duplicada para mesmo servico",
			clienteID: clienteID,
			limpezaID: limpezaID,
			nota:      5,
			wantErr:   true,
			wantErrIs: errosdominio.ErrAvaliacaoDuplicada,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := entity.NovaAvaliacao(tt.limpezaID, profissionalID, tt.clienteID, tt.nota, "Ótimo serviço")
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

func TestAvaliacao_ConsultarAvaliacaoQueClienteDeuParaServico(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax2@aval.com", "fax2aval")
	clienteID := inserirUsuario(t, db, "cli2@aval.com", "cli2aval")
	limpezaID := inserirLimpeza(t, db, profissionalID, "Servico2 Aval")

	a := entity.NovaAvaliacao(limpezaID, profissionalID, clienteID, 5, "Perfeito")
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
		{name: "avaliacao existente retorna nota e comentario", clienteID: clienteID, limpezaID: limpezaID},
		{
			name:      "avaliacao inexistente retorna erro de nao encontrada",
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

func TestAvaliacao_ListarTodasAvaliacoesRecebidasPeloProfissional(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax3@aval.com", "fax3aval")
	cli1 := inserirUsuario(t, db, "cli3a@aval.com", "cli3a")
	cli2 := inserirUsuario(t, db, "cli3b@aval.com", "cli3b")
	limp1 := inserirLimpeza(t, db, profissionalID, "Aval1")
	limp2 := inserirLimpeza(t, db, profissionalID, "Aval2")

	repo.Salvar(ctx, entity.NovaAvaliacao(limp1, profissionalID, cli1, 4, "Bom"))
	repo.Salvar(ctx, entity.NovaAvaliacao(limp2, profissionalID, cli2, 5, "Excelente"))

	got, err := repo.ListarPorProfissional(ctx, profissionalID)
	if err != nil {
		t.Fatalf("ListarPorProfissional() error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d; want 2", len(got))
	}
}

func TestAvaliacao_CalcularMediaEQuantidadeDeAvaliacoes(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAvaliacaoPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax4@aval.com", "fax4aval")
	cli1 := inserirUsuario(t, db, "cli4a@aval.com", "cli4a")
	cli2 := inserirUsuario(t, db, "cli4b@aval.com", "cli4b")
	limp1 := inserirLimpeza(t, db, profissionalID, "Aval3")
	limp2 := inserirLimpeza(t, db, profissionalID, "Aval4")

	repo.Salvar(ctx, entity.NovaAvaliacao(limp1, profissionalID, cli1, 4, ""))
	repo.Salvar(ctx, entity.NovaAvaliacao(limp2, profissionalID, cli2, 2, ""))

	t.Run("media calculada corretamente com duas avaliacoes", func(t *testing.T) {
		ag, err := repo.BuscarAgregadoPorProfissional(ctx, profissionalID)
		if err != nil {
			t.Fatalf("BuscarAgregadoPorProfissional() error: %v", err)
		}
		if ag.TotalAvaliacoes != 2 {
			t.Errorf("TotalAvaliacoes = %d; want 2", ag.TotalAvaliacoes)
		}
		if ag.MediaNota != 3.0 {
			t.Errorf("MediaNota = %v; want 3.0", ag.MediaNota)
		}
	})

	t.Run("profissional sem avaliacoes retorna media zero", func(t *testing.T) {
		semAvaliacaoID := inserirUsuario(t, db, "semav@aval.com", "semavaliacao")
		ag, err := repo.BuscarAgregadoPorProfissional(ctx, semAvaliacaoID)
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
