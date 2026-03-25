//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
	"limpaGo/infra/postgres"
)

func TestRepositorioSolicitacaoPG_Salvar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioSolicitacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "faxs@sol.com", "faxsol")
	clienteID := inserirUsuario(t, db, "clisol@sol.com", "clisol")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico Sol")

	s := &entity.Solicitacao{
		ClienteID:    clienteID,
		LimpezaID:    limpezaID,
		Status:       valueobject.StatusSolicitacaoPendente,
		DataAgendada: time.Now().Add(48 * time.Hour),
		PrecoTotal:   150.0,
		Endereco: valueobject.Endereco{
			Rua:    "Av. Paulista",
			Bairro: "Bela Vista",
			Cidade: "São Paulo",
			Estado: "SP",
			CEP:    "01311-200",
		},
	}

	if err := repo.Salvar(ctx, s); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}
	if s.ID == 0 {
		t.Error("ID = 0; want > 0")
	}
	if s.CriadoEm.IsZero() {
		t.Error("CriadoEm zerado; want preenchido")
	}
}

func TestRepositorioSolicitacaoPG_BuscarPorClienteELimpeza(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioSolicitacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax2s@sol.com", "fax2sol")
	clienteID := inserirUsuario(t, db, "cli2sol@sol.com", "cli2sol")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico2 Sol")
	inserirSolicitacao(t, db, clienteID, limpezaID)

	tests := []struct {
		name      string
		clienteID int
		limpezaID int
		wantErr   error
	}{
		{name: "encontrada", clienteID: clienteID, limpezaID: limpezaID},
		{name: "nao encontrada", clienteID: 999999, limpezaID: limpezaID, wantErr: errosdominio.ErrSolicitacaoNaoEncontrada},
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
				t.Fatalf("BuscarPorClienteELimpeza() error: %v", err)
			}
			if got == nil {
				t.Error("got nil; want solicitacao")
			}
		})
	}
}

func TestRepositorioSolicitacaoPG_BuscarAtivaPorClienteELimpeza(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioSolicitacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax3s@sol.com", "fax3sol")
	clienteID := inserirUsuario(t, db, "cli3sol@sol.com", "cli3sol")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico3 Sol")

	// Inserir solicitação pendente
	solID := inserirSolicitacao(t, db, clienteID, limpezaID)

	t.Run("pendente retorna ativa", func(t *testing.T) {
		got, err := repo.BuscarAtivaPorClienteELimpeza(ctx, clienteID, limpezaID)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil; want ativa")
		}
	})

	// Cancelar a solicitação
	db.Exec(`UPDATE solicitacoes SET status='cancelada' WHERE id=$1`, solID)

	t.Run("cancelada nao retorna ativa", func(t *testing.T) {
		_, err := repo.BuscarAtivaPorClienteELimpeza(ctx, clienteID, limpezaID)
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada) {
			t.Errorf("got %v; want %v", err, errosdominio.ErrSolicitacaoNaoEncontrada)
		}
	})
}

func TestRepositorioSolicitacaoPG_Atualizar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioSolicitacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax4s@sol.com", "fax4sol")
	clienteID := inserirUsuario(t, db, "cli4sol@sol.com", "cli4sol")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico4 Sol")
	solID := inserirSolicitacao(t, db, clienteID, limpezaID)

	s := &entity.Solicitacao{
		ID:     solID,
		Status: valueobject.StatusSolicitacaoAceita,
	}
	if err := repo.Atualizar(ctx, s); err != nil {
		t.Fatalf("Atualizar() error: %v", err)
	}
	if s.AtualizadoEm.IsZero() {
		t.Error("AtualizadoEm zerado; want preenchido")
	}
}

func TestRepositorioSolicitacaoPG_Deletar(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioSolicitacaoPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax5s@sol.com", "fax5sol")
	clienteID := inserirUsuario(t, db, "cli5sol@sol.com", "cli5sol")
	limpezaID := inserirLimpeza(t, db, faxineiroID, "Servico5 Sol")
	inserirSolicitacao(t, db, clienteID, limpezaID)

	if err := repo.Deletar(ctx, clienteID, limpezaID); err != nil {
		t.Fatalf("Deletar() error: %v", err)
	}

	_, err := repo.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if !errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada) {
		t.Errorf("apos deletar: got %v; want %v", err, errosdominio.ErrSolicitacaoNaoEncontrada)
	}
}
