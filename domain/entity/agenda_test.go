package entity

import (
	"errors"
	"testing"
	"time"

	errosdominio "limpaGo/domain/errors"
)

func TestNovaDisponibilidade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		faxID      int
		dia        time.Weekday
		horaInicio int
		horaFim    int
		wantErr    bool
		wantCampo  string
	}{
		{"valores validos", 1, time.Monday, 8, 12, false, ""},
		{"limites extremos 0 a 24", 1, time.Sunday, 0, 24, false, ""},
		{"hora inicio negativa", 1, time.Monday, -1, 12, true, "hora_inicio"},
		{"hora inicio 24", 1, time.Monday, 24, 25, true, "hora_inicio"},
		{"hora fim zero", 1, time.Monday, 0, 0, true, "hora_fim"},
		{"hora fim 25", 1, time.Monday, 8, 25, true, "hora_fim"},
		{"hora fim igual inicio", 1, time.Monday, 10, 10, true, "hora_fim"},
		{"hora fim menor que inicio", 1, time.Monday, 12, 8, true, "hora_fim"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d, err := NovaDisponibilidade(tt.faxID, tt.dia, tt.horaInicio, tt.horaFim)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NovaDisponibilidade() error = nil; want error")
				}
				var ev *ErroValidacao
				if tt.wantCampo != "" && errors.As(err, &ev) && ev.Campo != tt.wantCampo {
					t.Errorf("ErroValidacao.Campo = %q; want %q", ev.Campo, tt.wantCampo)
				}
				return
			}

			if err != nil {
				t.Fatalf("NovaDisponibilidade() unexpected error: %v", err)
			}
			if d.ProfissionalID != tt.faxID {
				t.Errorf("ProfissionalID = %d; want %d", d.ProfissionalID, tt.faxID)
			}
			if d.HoraInicio != tt.horaInicio {
				t.Errorf("HoraInicio = %d; want %d", d.HoraInicio, tt.horaInicio)
			}
			if d.HoraFim != tt.horaFim {
				t.Errorf("HoraFim = %d; want %d", d.HoraFim, tt.horaFim)
			}
		})
	}
}

func TestDisponibilidade_DuracaoHoras(t *testing.T) {
	t.Parallel()

	d := &Disponibilidade{HoraInicio: 8, HoraFim: 12}
	got := d.DuracaoHoras()

	if got != 4 {
		t.Errorf("DuracaoHoras() = %d; want 4", got)
	}
}

func TestNovoBloqueioServico(t *testing.T) {
	t.Parallel()

	futuro := time.Now().Add(48 * time.Hour)
	futuroFim := futuro.Add(3 * time.Hour)
	passado := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name      string
		inicio    time.Time
		fim       time.Time
		wantErr   bool
		wantErrIs error
	}{
		{"datas futuras validas", futuro, futuroFim, false, nil},
		{"data fim antes de inicio", futuroFim, futuro, true, nil},
		{"data fim igual inicio", futuro, futuro, true, nil},
		{"data inicio no passado", passado, futuro, true, errosdominio.ErrAgendamentoNoPassado},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b, err := NovoBloqueioServico(1, 100, tt.inicio, tt.fim)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NovoBloqueioServico() error = nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("error = %v; want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("NovoBloqueioServico() unexpected error: %v", err)
			}
			if b.SolicitacaoID == nil || *b.SolicitacaoID != 100 {
				t.Error("SolicitacaoID should be 100")
			}
		})
	}
}

func TestNovoBloqueiopessoal(t *testing.T) {
	t.Parallel()

	futuro := time.Now().Add(48 * time.Hour)
	futuroFim := futuro.Add(3 * time.Hour)

	b, err := NovoBloqueiopessoal(1, futuro, futuroFim)
	if err != nil {
		t.Fatalf("NovoBloqueiopessoal() unexpected error: %v", err)
	}
	if b.SolicitacaoID != nil {
		t.Error("SolicitacaoID = non-nil; want nil for bloqueio pessoal")
	}
}

func TestBloqueio_EPessoal(t *testing.T) {
	t.Parallel()

	solicitacaoID := 1
	tests := []struct {
		name          string
		solicitacaoID *int
		want          bool
	}{
		{"bloqueio pessoal", nil, true},
		{"bloqueio de servico", &solicitacaoID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &Bloqueio{SolicitacaoID: tt.solicitacaoID}
			got := b.EPessoal()
			if got != tt.want {
				t.Errorf("EPessoal() = %v; want %v", got, tt.want)
			}
		})
	}
}
