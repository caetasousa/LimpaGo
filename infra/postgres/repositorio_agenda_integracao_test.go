//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/infra/postgres"
)

func TestAgenda_FaxineiroDefineHorariosDisponiveisNaSemana(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAgendaPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax@agenda.com", "faxagenda")

	var dispID int

	t.Run("faxineiro cadastra disponibilidade na segunda-feira", func(t *testing.T) {
		d := &entity.Disponibilidade{
			FaxineiroID: faxineiroID,
			DiaSemana:   time.Monday,
			HoraInicio:  8,
			HoraFim:     17,
		}
		if err := repo.SalvarDisponibilidade(ctx, d); err != nil {
			t.Fatalf("SalvarDisponibilidade() error: %v", err)
		}
		if d.ID == 0 {
			t.Error("ID = 0; want > 0")
		}
		dispID = d.ID
	})

	t.Run("consulta retorna todas as disponibilidades do faxineiro", func(t *testing.T) {
		got, err := repo.ListarDisponibilidadePorFaxineiro(ctx, faxineiroID)
		if err != nil {
			t.Fatalf("ListarDisponibilidadePorFaxineiro() error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("len = %d; want 1", len(got))
		}
	})

	t.Run("filtrar disponibilidade por dia retorna apenas segunda-feira", func(t *testing.T) {
		got, err := repo.ListarDisponibilidadePorDia(ctx, faxineiroID, time.Monday)
		if err != nil {
			t.Fatalf("ListarDisponibilidadePorDia() error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("len = %d; want 1", len(got))
		}
	})

	t.Run("dia sem disponibilidade retorna lista vazia", func(t *testing.T) {
		got, err := repo.ListarDisponibilidadePorDia(ctx, faxineiroID, time.Sunday)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("len = %d; want 0", len(got))
		}
	})

	t.Run("faxineiro remove disponibilidade da agenda", func(t *testing.T) {
		if err := repo.DeletarDisponibilidade(ctx, dispID, faxineiroID); err != nil {
			t.Fatalf("DeletarDisponibilidade() error: %v", err)
		}
		got, _ := repo.ListarDisponibilidadePorFaxineiro(ctx, faxineiroID)
		if len(got) != 0 {
			t.Errorf("apos deletar len = %d; want 0", len(got))
		}
	})
}

func TestAgenda_FaxineiroBloqueiaHorariosParaCompromissosPessoais(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioAgendaPG(db)
	ctx := context.Background()

	faxineiroID := inserirUsuario(t, db, "fax2@agenda.com", "fax2agenda")

	agora := time.Now().Truncate(time.Second).UTC()
	inicio := agora.Add(24 * time.Hour)
	fim := inicio.Add(4 * time.Hour)

	var bloqueioID int

	t.Run("faxineiro cria bloqueio pessoal na agenda", func(t *testing.T) {
		b := &entity.Bloqueio{
			FaxineiroID: faxineiroID,
			DataInicio:  inicio,
			DataFim:     fim,
		}
		if err := repo.SalvarBloqueio(ctx, b); err != nil {
			t.Fatalf("SalvarBloqueio() error: %v", err)
		}
		if b.ID == 0 {
			t.Error("ID = 0; want > 0")
		}
		bloqueioID = b.ID
	})

	t.Run("bloqueio pode ser consultado pelo ID", func(t *testing.T) {
		got, err := repo.BuscarBloqueioPorID(ctx, bloqueioID)
		if err != nil {
			t.Fatalf("BuscarBloqueioPorID() error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil; want bloqueio")
		}
		if got.FaxineiroID != faxineiroID {
			t.Errorf("FaxineiroID = %d; want %d", got.FaxineiroID, faxineiroID)
		}
	})

	t.Run("listar bloqueios retorna todos do faxineiro", func(t *testing.T) {
		got, err := repo.ListarBloqueiosPorFaxineiro(ctx, faxineiroID)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("len = %d; want 1", len(got))
		}
	})

	t.Run("sistema detecta conflito quando periodo se sobrepoe ao bloqueio", func(t *testing.T) {
		// Período que se sobrepõe ao bloqueio
		got, err := repo.ListarBloqueiosPorPeriodo(ctx, faxineiroID,
			inicio.Add(-1*time.Hour), fim.Add(1*time.Hour))
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("len = %d; want 1 (overlap detectado)", len(got))
		}
	})

	t.Run("periodo sem sobreposicao nao retorna bloqueios", func(t *testing.T) {
		// Período que não se sobrepõe
		got, err := repo.ListarBloqueiosPorPeriodo(ctx, faxineiroID,
			fim.Add(2*time.Hour), fim.Add(6*time.Hour))
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("len = %d; want 0 (sem overlap)", len(got))
		}
	})

	t.Run("bloqueio inexistente retorna erro de nao encontrado", func(t *testing.T) {
		_, err := repo.BuscarBloqueioPorID(ctx, 999999)
		if !errors.Is(err, errosdominio.ErrBloqueioNaoEncontrado) {
			t.Errorf("got %v; want %v", err, errosdominio.ErrBloqueioNaoEncontrado)
		}
	})

	t.Run("faxineiro remove bloqueio da agenda", func(t *testing.T) {
		if err := repo.DeletarBloqueio(ctx, bloqueioID); err != nil {
			t.Fatalf("DeletarBloqueio() error: %v", err)
		}
		_, err := repo.BuscarBloqueioPorID(ctx, bloqueioID)
		if !errors.Is(err, errosdominio.ErrBloqueioNaoEncontrado) {
			t.Errorf("apos deletar: got %v; want ErrBloqueioNaoEncontrado", err)
		}
	})
}
