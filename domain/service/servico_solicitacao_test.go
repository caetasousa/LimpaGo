package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
	"limpaGo/domain/valueobject"
)

// setupSolicitacao cria cenario completo: faxineiroID=1, limpeza 3h, disponibilidade segunda 8-17h.
func setupSolicitacao(t *testing.T) (
	*service.ServicoSolicitacao,
	*service.ServicoAgenda,
	*testutil.RepositorioLimpezaMock,
	*entity.Limpeza,
	time.Time,
) {
	t.Helper()

	repoLimpeza := testutil.NovoRepositorioLimpezaMock()
	repoSolic := testutil.NovoRepositorioSolicitacaoMock()
	repoAgenda := testutil.NovoRepositorioAgendaMock()

	svcAgenda := service.NovoServicoAgenda(repoAgenda)
	svcSolic := service.NovoServicoSolicitacao(repoSolic, repoLimpeza, svcAgenda)

	ctx := context.Background()

	limpeza := &entity.Limpeza{
		FaxineiroID:     1,
		Nome:            "Limpeza Teste",
		ValorHora:       50,
		DuracaoEstimada: 3,
		TipoLimpeza:     valueobject.TipoLimpezaPadrao,
	}
	_ = repoLimpeza.Salvar(ctx, limpeza)

	dataAgendada := proximoDia(time.Monday, 10)
	_, _ = svcAgenda.AdicionarDisponibilidade(ctx, 1, time.Monday, 8, 17)

	return svcSolic, svcAgenda, repoLimpeza, limpeza, dataAgendada
}

func TestServicoSolicitacao_CriarSolicitacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("cria solicitacao pendente", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)

		s, err := svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)
		if err != nil {
			t.Fatalf("CriarSolicitacao() unexpected error: %v", err)
		}
		if s.Status != valueobject.StatusSolicitacaoPendente {
			t.Errorf("Status = %q; want pendente", s.Status)
		}
		if s.PrecoTotal != 150 {
			t.Errorf("PrecoTotal = %f; want 150", s.PrecoTotal)
		}
	})

	t.Run("duplicada ativa", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		_, err := svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)
		if !errors.Is(err, errosdominio.ErrSolicitacaoDuplicada) {
			t.Errorf("error = %v; want ErrSolicitacaoDuplicada", err)
		}
	})

	t.Run("faxineiro solicitando proprio servico", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)

		_, err := svc.CriarSolicitacao(ctx, 1, limpeza.ID, dataAg)
		if !errors.Is(err, errosdominio.ErrFaxineiroNaoPodeSolicitarProprio) {
			t.Errorf("error = %v; want ErrFaxineiroNaoPodeSolicitarProprio", err)
		}
	})

	t.Run("data no passado", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, _ := setupSolicitacao(t)
		passado := time.Now().Add(-24 * time.Hour)

		_, err := svc.CriarSolicitacao(ctx, 2, limpeza.ID, passado)
		if err == nil {
			t.Fatal("error = nil; want error for past date")
		}
	})

	t.Run("sem disponibilidade", func(t *testing.T) {
		t.Parallel()
		repoLimpeza := testutil.NovoRepositorioLimpezaMock()
		repoSolic := testutil.NovoRepositorioSolicitacaoMock()
		repoAgenda := testutil.NovoRepositorioAgendaMock()
		svcAgenda := service.NovoServicoAgenda(repoAgenda)
		svcSolic := service.NovoServicoSolicitacao(repoSolic, repoLimpeza, svcAgenda)

		limpeza := &entity.Limpeza{FaxineiroID: 1, ValorHora: 50, DuracaoEstimada: 3, TipoLimpeza: valueobject.TipoLimpezaPadrao}
		_ = repoLimpeza.Salvar(ctx, limpeza)

		futuro := time.Now().Add(48 * time.Hour)
		_, err := svcSolic.CriarSolicitacao(ctx, 2, limpeza.ID, futuro)
		if !errors.Is(err, errosdominio.ErrHorarioIndisponivel) {
			t.Errorf("error = %v; want ErrHorarioIndisponivel", err)
		}
	})
}

func TestServicoSolicitacao_AceitarSolicitacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("aceita e cria bloqueio", func(t *testing.T) {
		t.Parallel()
		svc, svcAgenda, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		s, err := svc.AceitarSolicitacao(ctx, 1, 2, limpeza.ID)
		if err != nil {
			t.Fatalf("AceitarSolicitacao() unexpected error: %v", err)
		}
		if s.Status != valueobject.StatusSolicitacaoAceita {
			t.Errorf("Status = %q; want aceita", s.Status)
		}

		bloqueios, _ := svcAgenda.ListarBloqueios(ctx, 1)
		if len(bloqueios) == 0 {
			t.Error("expected block created in agenda")
		}
	})

	t.Run("nao faxineiro", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		_, err := svc.AceitarSolicitacao(ctx, 999, 2, limpeza.ID)
		if !errors.Is(err, errosdominio.ErrNaoEFaxineiroDaSolicitacao) {
			t.Errorf("error = %v; want ErrNaoEFaxineiroDaSolicitacao", err)
		}
	})

	t.Run("solicitacao nao encontrada", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, _ := setupSolicitacao(t)

		_, err := svc.AceitarSolicitacao(ctx, 1, 999, limpeza.ID)
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada) {
			t.Errorf("error = %v; want ErrSolicitacaoNaoEncontrada", err)
		}
	})

	t.Run("ja aceita retorna erro", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)
		_, _ = svc.AceitarSolicitacao(ctx, 1, 2, limpeza.ID)

		_, err := svc.AceitarSolicitacao(ctx, 1, 2, limpeza.ID)
		if err == nil {
			t.Fatal("error = nil; want error for already accepted")
		}
	})
}

func TestServicoSolicitacao_RejeitarSolicitacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("rejeita pendente", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		s, err := svc.RejeitarSolicitacao(ctx, 1, 2, limpeza.ID)
		if err != nil {
			t.Fatalf("RejeitarSolicitacao() unexpected error: %v", err)
		}
		if s.Status != valueobject.StatusSolicitacaoRejeitada {
			t.Errorf("Status = %q; want rejeitada", s.Status)
		}
	})

	t.Run("nao faxineiro", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		_, err := svc.RejeitarSolicitacao(ctx, 999, 2, limpeza.ID)
		if !errors.Is(err, errosdominio.ErrNaoEFaxineiroDaSolicitacao) {
			t.Errorf("error = %v; want ErrNaoEFaxineiroDaSolicitacao", err)
		}
	})

	t.Run("rejeitar aceita", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)
		_, _ = svc.AceitarSolicitacao(ctx, 1, 2, limpeza.ID)

		_, err := svc.RejeitarSolicitacao(ctx, 1, 2, limpeza.ID)
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerRejeitada) {
			t.Errorf("error = %v; want ErrSolicitacaoNaoPodeSerRejeitada", err)
		}
	})
}

func TestServicoSolicitacao_CancelarSolicitacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("cancela pendente sem multa", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

		s, err := svc.CancelarSolicitacao(ctx, 2, limpeza.ID)
		if err != nil {
			t.Fatalf("CancelarSolicitacao() unexpected error: %v", err)
		}
		if s.Status != valueobject.StatusSolicitacaoCancelada {
			t.Errorf("Status = %q; want cancelada", s.Status)
		}
		if s.MultaCancelamento != 0 {
			t.Errorf("MultaCancelamento = %f; want 0", s.MultaCancelamento)
		}
	})

	t.Run("cancela aceita libera bloqueio", func(t *testing.T) {
		t.Parallel()
		svc, svcAgenda, _, limpeza, dataAg := setupSolicitacao(t)
		_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)
		_, _ = svc.AceitarSolicitacao(ctx, 1, 2, limpeza.ID)

		s, err := svc.CancelarSolicitacao(ctx, 2, limpeza.ID)
		if err != nil {
			t.Fatalf("CancelarSolicitacao() unexpected error: %v", err)
		}
		if s.Status != valueobject.StatusSolicitacaoCancelada {
			t.Errorf("Status = %q; want cancelada", s.Status)
		}

		bloqueios, _ := svcAgenda.ListarBloqueios(ctx, 1)
		if len(bloqueios) != 0 {
			t.Errorf("len(bloqueios) = %d; want 0 after cancel", len(bloqueios))
		}
	})

	t.Run("nao encontrada", func(t *testing.T) {
		t.Parallel()
		svc, _, _, limpeza, _ := setupSolicitacao(t)

		_, err := svc.CancelarSolicitacao(ctx, 999, limpeza.ID)
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada) {
			t.Errorf("error = %v; want ErrSolicitacaoNaoEncontrada", err)
		}
	})
}

func TestServicoSolicitacao_ListarSolicitacoesPorLimpeza(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc, _, _, limpeza, dataAg := setupSolicitacao(t)
	_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

	lista, err := svc.ListarSolicitacoesPorLimpeza(ctx, 1, limpeza.ID)
	if err != nil {
		t.Fatalf("ListarSolicitacoesPorLimpeza() unexpected error: %v", err)
	}
	if len(lista) != 1 {
		t.Errorf("len(lista) = %d; want 1", len(lista))
	}
}

func TestServicoSolicitacao_ListarSolicitacoesPorCliente(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc, _, _, limpeza, dataAg := setupSolicitacao(t)
	_, _ = svc.CriarSolicitacao(ctx, 2, limpeza.ID, dataAg)

	lista, err := svc.ListarSolicitacoesPorCliente(ctx, 2)
	if err != nil {
		t.Fatalf("ListarSolicitacoesPorCliente() unexpected error: %v", err)
	}
	if len(lista) != 1 {
		t.Errorf("len(lista) = %d; want 1", len(lista))
	}
}
