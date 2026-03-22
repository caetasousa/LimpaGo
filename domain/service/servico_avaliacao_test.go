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

func setupServicoAvaliacao(t *testing.T) (
	*service.ServicoAvaliacao,
	*testutil.RepositorioAvaliacaoMock,
	*testutil.RepositorioSolicitacaoMock,
	*testutil.RepositorioLimpezaMock,
) {
	t.Helper()

	repoAval := testutil.NovoRepositorioAvaliacaoMock()
	repoSolic := testutil.NovoRepositorioSolicitacaoMock()
	repoLimpeza := testutil.NovoRepositorioLimpezaMock()

	svc := service.NovoServicoAvaliacao(repoAval, repoSolic, repoLimpeza)

	ctx := context.Background()

	limpeza := &entity.Limpeza{
		FaxineiroID:     1,
		Nome:            "Limpeza Teste",
		ValorHora:       50,
		DuracaoEstimada: 3,
		TipoLimpeza:     valueobject.TipoLimpezaPadrao,
	}
	_ = repoLimpeza.Salvar(ctx, limpeza)

	solicitacao := &entity.Solicitacao{
		ClienteID:    2,
		LimpezaID:    limpeza.ID,
		Status:       valueobject.StatusSolicitacaoAceita,
		DataAgendada: time.Now().Add(48 * time.Hour),
		PrecoTotal:   150,
	}
	_ = repoSolic.Salvar(ctx, solicitacao)

	return svc, repoAval, repoSolic, repoLimpeza
}

func TestServicoAvaliacao_CriarAvaliacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("avaliacao com nota e comentario", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _ := setupServicoAvaliacao(t)

		a, err := svc.CriarAvaliacao(ctx, 2, 1, 5, "Excelente!")
		if err != nil {
			t.Fatalf("CriarAvaliacao() unexpected error: %v", err)
		}
		if int(a.Nota) != 5 {
			t.Errorf("Nota = %d; want 5", a.Nota)
		}
		if a.FaxineiroID != 1 {
			t.Errorf("FaxineiroID = %d; want 1", a.FaxineiroID)
		}
	})

	t.Run("nota zero valida", func(t *testing.T) {
		t.Parallel()
		svc, _, repoSolic, repoLimpeza := setupServicoAvaliacao(t)

		limpeza2 := &entity.Limpeza{FaxineiroID: 1, ValorHora: 40, DuracaoEstimada: 2, TipoLimpeza: valueobject.TipoLimpezaPesada}
		_ = repoLimpeza.Salvar(ctx, limpeza2)
		solic2 := &entity.Solicitacao{ClienteID: 2, LimpezaID: limpeza2.ID, Status: valueobject.StatusSolicitacaoAceita, PrecoTotal: 80}
		_ = repoSolic.Salvar(ctx, solic2)

		a, err := svc.CriarAvaliacao(ctx, 2, limpeza2.ID, 0, "")
		if err != nil {
			t.Fatalf("CriarAvaliacao() unexpected error: %v", err)
		}
		if int(a.Nota) != 0 {
			t.Errorf("Nota = %d; want 0", a.Nota)
		}
	})

	t.Run("solicitacao inexistente", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _ := setupServicoAvaliacao(t)

		_, err := svc.CriarAvaliacao(ctx, 999, 999, 5, "")
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada) {
			t.Errorf("error = %v; want ErrSolicitacaoNaoEncontrada", err)
		}
	})

	t.Run("solicitacao pendente", func(t *testing.T) {
		t.Parallel()
		svc, _, repoSolic, repoLimpeza := setupServicoAvaliacao(t)

		limpeza3 := &entity.Limpeza{FaxineiroID: 1, ValorHora: 30, DuracaoEstimada: 1, TipoLimpeza: valueobject.TipoLimpezaExpress}
		_ = repoLimpeza.Salvar(ctx, limpeza3)
		solicPend := &entity.Solicitacao{ClienteID: 3, LimpezaID: limpeza3.ID, Status: valueobject.StatusSolicitacaoPendente, PrecoTotal: 30}
		_ = repoSolic.Salvar(ctx, solicPend)

		_, err := svc.CriarAvaliacao(ctx, 3, limpeza3.ID, 5, "")
		if !errors.Is(err, errosdominio.ErrSolicitacaoNaoAceita) {
			t.Errorf("error = %v; want ErrSolicitacaoNaoAceita", err)
		}
	})

	t.Run("avaliacao duplicada", func(t *testing.T) {
		t.Parallel()
		svc, repoAval, repoSolic, repoLimpeza := setupServicoAvaliacao(t)

		limpezaDup := &entity.Limpeza{FaxineiroID: 1, ValorHora: 30, DuracaoEstimada: 1, TipoLimpeza: valueobject.TipoLimpezaExpress}
		_ = repoLimpeza.Salvar(ctx, limpezaDup)
		solicDup := &entity.Solicitacao{ClienteID: 6, LimpezaID: limpezaDup.ID, Status: valueobject.StatusSolicitacaoAceita, PrecoTotal: 30}
		_ = repoSolic.Salvar(ctx, solicDup)

		nota, _ := valueobject.NovaNota(5)
		avalExistente := entity.NovaAvaliacao(limpezaDup.ID, 1, 6, nota, "Já existe")
		_ = repoAval.Salvar(ctx, avalExistente)

		_, err := svc.CriarAvaliacao(ctx, 6, limpezaDup.ID, 4, "Outra")
		if !errors.Is(err, errosdominio.ErrAvaliacaoDuplicada) {
			t.Errorf("error = %v; want ErrAvaliacaoDuplicada", err)
		}
	})

	notasInvalidas := []struct {
		name string
		nota int
	}{
		{"nota 6 invalida", 6},
		{"nota -1 invalida", -1},
	}

	for _, tt := range notasInvalidas {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, _, repoSolic, repoLimpeza := setupServicoAvaliacao(t)

			l := &entity.Limpeza{FaxineiroID: 1, ValorHora: 30, DuracaoEstimada: 1, TipoLimpeza: valueobject.TipoLimpezaExpress}
			_ = repoLimpeza.Salvar(ctx, l)
			s := &entity.Solicitacao{ClienteID: 10, LimpezaID: l.ID, Status: valueobject.StatusSolicitacaoAceita, PrecoTotal: 30}
			_ = repoSolic.Salvar(ctx, s)

			_, err := svc.CriarAvaliacao(ctx, 10, l.ID, tt.nota, "")
			var ev *entity.ErroValidacao
			if !errors.As(err, &ev) || ev.Campo != "nota" {
				t.Errorf("error = %v; want ErroValidacao{Campo: nota}", err)
			}
		})
	}
}

func TestServicoAvaliacao_BuscarEstatisticasFaxineiro(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc, _, _, _ := setupServicoAvaliacao(t)
	_, _ = svc.CriarAvaliacao(ctx, 2, 1, 5, "")

	agg, err := svc.BuscarEstatisticasFaxineiro(ctx, 1)
	if err != nil {
		t.Fatalf("BuscarEstatisticasFaxineiro() unexpected error: %v", err)
	}
	if agg.TotalAvaliacoes != 1 {
		t.Errorf("TotalAvaliacoes = %d; want 1", agg.TotalAvaliacoes)
	}
	if agg.MediaNota != 5 {
		t.Errorf("MediaNota = %f; want 5", agg.MediaNota)
	}
}

func TestServicoAvaliacao_ListarAvaliacoesPorFaxineiro(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc, _, _, _ := setupServicoAvaliacao(t)
	_, _ = svc.CriarAvaliacao(ctx, 2, 1, 4, "Bom")

	lista, err := svc.ListarAvaliacoesPorFaxineiro(ctx, 1)
	if err != nil {
		t.Fatalf("ListarAvaliacoesPorFaxineiro() unexpected error: %v", err)
	}
	if len(lista) != 1 {
		t.Errorf("len(lista) = %d; want 1", len(lista))
	}
}
