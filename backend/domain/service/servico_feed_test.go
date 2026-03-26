package service_test

import (
	"context"
	"testing"
	"time"

	"limpaGo/domain/entity"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
	"limpaGo/domain/valueobject"
)

func setupServicoFeed(t *testing.T) (*service.ServicoFeed, *testutil.RepositorioFeedMock) {
	t.Helper()
	repo := testutil.NovoRepositorioFeedMock()
	svc := service.NovoServicoFeed(repo)
	return svc, repo
}

func TestServicoFeed_BuscarFeed(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("retorna pagina com itens", func(t *testing.T) {
		t.Parallel()
		svc, repo := setupServicoFeed(t)
		repo.Itens = []*entity.ItemFeed{
			{Limpeza: &entity.Limpeza{Nome: "L1"}, TipoEvento: valueobject.TipoEventoFeedCriacao, DataEvento: time.Now()},
			{Limpeza: &entity.Limpeza{Nome: "L2"}, TipoEvento: valueobject.TipoEventoFeedAtualizacao, DataEvento: time.Now()},
		}

		pf, err := svc.BuscarFeed(ctx, 1, 10)
		if err != nil {
			t.Fatalf("BuscarFeed() unexpected error: %v", err)
		}
		if len(pf.Itens) != 2 {
			t.Errorf("len(Itens) = %d; want 2", len(pf.Itens))
		}
		if pf.TotalItens != 2 {
			t.Errorf("TotalItens = %d; want 2", pf.TotalItens)
		}
	})

	t.Run("paginacao respeitada", func(t *testing.T) {
		t.Parallel()
		svc, repo := setupServicoFeed(t)
		for i := 0; i < 5; i++ {
			repo.Itens = append(repo.Itens, &entity.ItemFeed{
				Limpeza:    &entity.Limpeza{Nome: "L"},
				TipoEvento: valueobject.TipoEventoFeedCriacao,
				DataEvento: time.Now(),
			})
		}

		pf, err := svc.BuscarFeed(ctx, 1, 3)
		if err != nil {
			t.Fatalf("BuscarFeed() unexpected error: %v", err)
		}
		if len(pf.Itens) != 3 {
			t.Errorf("len(Itens) = %d; want 3", len(pf.Itens))
		}
	})

	t.Run("paginacao invalida corrigida", func(t *testing.T) {
		t.Parallel()
		svc, repo := setupServicoFeed(t)
		repo.Itens = []*entity.ItemFeed{
			{Limpeza: &entity.Limpeza{Nome: "L1"}, TipoEvento: valueobject.TipoEventoFeedCriacao, DataEvento: time.Now()},
		}

		pf, err := svc.BuscarFeed(ctx, 0, 0)
		if err != nil {
			t.Fatalf("BuscarFeed() unexpected error: %v", err)
		}
		if pf.Pagina != 1 {
			t.Errorf("Pagina = %d; want 1 (corrected)", pf.Pagina)
		}
	})

	t.Run("sem itens pagina vazia", func(t *testing.T) {
		t.Parallel()
		svc, _ := setupServicoFeed(t)

		pf, err := svc.BuscarFeed(ctx, 1, 10)
		if err != nil {
			t.Fatalf("BuscarFeed() unexpected error: %v", err)
		}
		if len(pf.Itens) != 0 {
			t.Errorf("len(Itens) = %d; want 0", len(pf.Itens))
		}
		if pf.TotalItens != 0 {
			t.Errorf("TotalItens = %d; want 0", pf.TotalItens)
		}
	})
}
