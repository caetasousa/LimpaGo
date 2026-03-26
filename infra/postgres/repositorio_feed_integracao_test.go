//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"limpaGo/infra/postgres"
)

func TestFeed_ClienteNavegaPelosServicosDisponiveisComPaginacao(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioFeedPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax@feed.com", "faxfeed")
	for i := 0; i < 5; i++ {
		inserirLimpeza(t, db, profissionalID, "Servico Feed")
	}

	tests := []struct {
		name          string
		pagina        int
		tamanhoPagina int
		wantLen       int
		wantTotal     int
	}{
		{
			name:          "primeira pagina mostra 3 de 5 servicos disponiveis",
			pagina:        1,
			tamanhoPagina: 3,
			wantLen:       3,
			wantTotal:     5,
		},
		{
			name:          "segunda pagina mostra os 2 servicos restantes",
			pagina:        2,
			tamanhoPagina: 3,
			wantLen:       2,
			wantTotal:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagina, err := repo.BuscarPaginaFeed(ctx, tt.pagina, tt.tamanhoPagina)
			if err != nil {
				t.Fatalf("BuscarPaginaFeed() error: %v", err)
			}
			if pagina == nil {
				t.Fatal("got nil; want pagina")
			}
			if len(pagina.Itens) != tt.wantLen {
				t.Errorf("Itens len = %d; want %d", len(pagina.Itens), tt.wantLen)
			}
			if pagina.TotalItens != tt.wantTotal {
				t.Errorf("TotalItens = %d; want %d", pagina.TotalItens, tt.wantTotal)
			}
		})
	}

	t.Run("feed sem servicos cadastrados retorna pagina vazia", func(t *testing.T) {
		limparTabelas(t, db)
		pagina, err := repo.BuscarPaginaFeed(ctx, 1, 10)
		if err != nil {
			t.Fatalf("BuscarPaginaFeed() error: %v", err)
		}
		if pagina == nil {
			t.Fatal("got nil; want pagina")
		}
		if len(pagina.Itens) != 0 {
			t.Errorf("Itens len = %d; want 0", len(pagina.Itens))
		}
	})
}
