package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func setupServicoAgenda(t *testing.T) *service.ServicoAgenda {
	t.Helper()
	return service.NovoServicoAgenda(testutil.NovoRepositorioAgendaMock())
}

// proximoDia retorna a proxima data futura que cai no dia da semana desejado.
func proximoDia(dia time.Weekday, hora int) time.Time {
	agora := time.Now()
	diff := int(dia) - int(agora.Weekday())
	if diff <= 0 {
		diff += 7
	}
	data := agora.AddDate(0, 0, diff)
	return time.Date(data.Year(), data.Month(), data.Day(), hora, 0, 0, 0, data.Location())
}

func TestServicoAgenda_AdicionarDisponibilidade(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name       string
		horaInicio int
		horaFim    int
		wantErr    bool
	}{
		{"horas validas", 8, 17, false},
		{"hora fim menor que inicio", 18, 8, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := setupServicoAgenda(t)

			d, err := svc.AdicionarDisponibilidade(ctx, 1, time.Monday, tt.horaInicio, tt.horaFim)

			if tt.wantErr {
				if err == nil {
					t.Fatal("AdicionarDisponibilidade() error = nil; want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("AdicionarDisponibilidade() unexpected error: %v", err)
			}
			if d.ID == 0 {
				t.Error("ID = 0; want assigned ID")
			}
		})
	}
}

func TestServicoAgenda_VerificarDisponibilidade(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("disponivel sem conflitos", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)
		_, _ = svc.AdicionarDisponibilidade(ctx, 1, time.Monday, 8, 17)

		inicio := proximoDia(time.Monday, 9)
		fim := proximoDia(time.Monday, 12)

		if err := svc.VerificarDisponibilidade(ctx, 1, inicio, fim); err != nil {
			t.Errorf("VerificarDisponibilidade() error = %v; want nil", err)
		}
	})

	t.Run("sem disponibilidade no dia", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)

		inicio := proximoDia(time.Monday, 9)
		fim := proximoDia(time.Monday, 12)

		err := svc.VerificarDisponibilidade(ctx, 1, inicio, fim)
		if !errors.Is(err, errosdominio.ErrHorarioIndisponivel) {
			t.Errorf("error = %v; want ErrHorarioIndisponivel", err)
		}
	})

	t.Run("disponibilidade insuficiente", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)
		_, _ = svc.AdicionarDisponibilidade(ctx, 1, time.Monday, 8, 10)

		inicio := proximoDia(time.Monday, 9)
		fim := proximoDia(time.Monday, 12)

		err := svc.VerificarDisponibilidade(ctx, 1, inicio, fim)
		if !errors.Is(err, errosdominio.ErrHorarioIndisponivel) {
			t.Errorf("error = %v; want ErrHorarioIndisponivel", err)
		}
	})

	t.Run("conflito com bloqueio existente", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)
		_, _ = svc.AdicionarDisponibilidade(ctx, 1, time.Monday, 8, 17)

		inicio := proximoDia(time.Monday, 9)
		fim := proximoDia(time.Monday, 12)
		_, _ = svc.CriarBloqueioPessoal(ctx, 1, inicio, fim)

		err := svc.VerificarDisponibilidade(ctx, 1, inicio, fim)
		if !errors.Is(err, errosdominio.ErrConflitoAgenda) {
			t.Errorf("error = %v; want ErrConflitoAgenda", err)
		}
	})
}

func TestServicoAgenda_CriarBloqueioServico(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc := setupServicoAgenda(t)
	futuro := time.Now().Add(48 * time.Hour)

	b, err := svc.CriarBloqueioServico(ctx, 1, 100, futuro, futuro.Add(3*time.Hour))
	if err != nil {
		t.Fatalf("CriarBloqueioServico() unexpected error: %v", err)
	}
	if b.SolicitacaoID == nil || *b.SolicitacaoID != 100 {
		t.Error("SolicitacaoID should be 100")
	}
	if b.ID == 0 {
		t.Error("ID = 0; want assigned ID")
	}
}

func TestServicoAgenda_CriarBloqueioPessoal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc := setupServicoAgenda(t)
	futuro := time.Now().Add(48 * time.Hour)

	b, err := svc.CriarBloqueioPessoal(ctx, 1, futuro, futuro.Add(3*time.Hour))
	if err != nil {
		t.Fatalf("CriarBloqueioPessoal() unexpected error: %v", err)
	}
	if !b.EPessoal() {
		t.Error("EPessoal() = false; want true")
	}
}

func TestServicoAgenda_RemoverBloqueioPessoal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(*service.ServicoAgenda) int
		faxID     int
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "remove pelo dono",
			setup: func(svc *service.ServicoAgenda) int {
				futuro := time.Now().Add(48 * time.Hour)
				b, _ := svc.CriarBloqueioPessoal(ctx, 1, futuro, futuro.Add(3*time.Hour))
				return b.ID
			},
			faxID: 1,
		},
		{
			name: "bloqueio inexistente",
			setup: func(_ *service.ServicoAgenda) int {
				return 999
			},
			faxID:     1,
			wantErr:   true,
			wantErrIs: errosdominio.ErrBloqueioNaoEncontrado,
		},
		{
			name: "outro profissional",
			setup: func(svc *service.ServicoAgenda) int {
				futuro := time.Now().Add(48 * time.Hour)
				b, _ := svc.CriarBloqueioPessoal(ctx, 1, futuro, futuro.Add(3*time.Hour))
				return b.ID
			},
			faxID:     999,
			wantErr:   true,
			wantErrIs: errosdominio.ErrNaoEProfissionalDoBloqueio,
		},
		{
			name: "bloqueio de servico",
			setup: func(svc *service.ServicoAgenda) int {
				futuro := time.Now().Add(48 * time.Hour)
				b, _ := svc.CriarBloqueioServico(ctx, 1, 100, futuro, futuro.Add(3*time.Hour))
				return b.ID
			},
			faxID:     1,
			wantErr:   true,
			wantErrIs: errosdominio.ErrBloqueioPessoalApenas,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := setupServicoAgenda(t)
			bloqueioID := tt.setup(svc)

			err := svc.RemoverBloqueioPessoal(ctx, bloqueioID, tt.faxID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("RemoverBloqueioPessoal() error = nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("error = %v; want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("RemoverBloqueioPessoal() unexpected error: %v", err)
			}
		})
	}
}

func TestServicoAgenda_ListarBloqueios(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc := setupServicoAgenda(t)
	futuro := time.Now().Add(48 * time.Hour)
	_, _ = svc.CriarBloqueioPessoal(ctx, 1, futuro, futuro.Add(2*time.Hour))
	_, _ = svc.CriarBloqueioServico(ctx, 1, 100, futuro.Add(3*time.Hour), futuro.Add(6*time.Hour))

	lista, err := svc.ListarBloqueios(ctx, 1)
	if err != nil {
		t.Fatalf("ListarBloqueios() unexpected error: %v", err)
	}
	if len(lista) != 2 {
		t.Errorf("len(lista) = %d; want 2", len(lista))
	}
}

func TestServicoAgenda_LiberarBloqueioPorSolicitacao(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("libera bloqueio existente", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)
		futuro := time.Now().Add(48 * time.Hour)
		_, _ = svc.CriarBloqueioServico(ctx, 1, 100, futuro, futuro.Add(3*time.Hour))

		if err := svc.LiberarBloqueioPorSolicitacao(ctx, 100); err != nil {
			t.Fatalf("LiberarBloqueioPorSolicitacao() unexpected error: %v", err)
		}

		lista, _ := svc.ListarBloqueios(ctx, 1)
		if len(lista) != 0 {
			t.Errorf("len(bloqueios) = %d; want 0 after release", len(lista))
		}
	})

	t.Run("sem bloqueio retorna nil", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoAgenda(t)

		if err := svc.LiberarBloqueioPorSolicitacao(ctx, 999); err != nil {
			t.Errorf("error = %v; want nil when no block exists", err)
		}
	})
}
