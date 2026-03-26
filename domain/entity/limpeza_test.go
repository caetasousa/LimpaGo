package entity

import (
	"errors"
	"testing"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

func TestNovaLimpeza(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		profissionalID     int
		nome            string
		valorHora       float64
		duracaoEstimada float64
		tipoLimpeza     valueobject.TipoLimpeza
		wantErr         bool
		wantCampo       string
	}{
		{
			name:            "parametros validos",
			profissionalID:     1,
			nome:            "Limpeza Padrão",
			valorHora:       50,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
		},
		{
			name:            "nome vazio",
			profissionalID:     1,
			nome:            "",
			valorHora:       50,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "nome",
		},
		{
			name:            "valor hora zero",
			profissionalID:     1,
			nome:            "Limpeza",
			valorHora:       0,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "valor_hora",
		},
		{
			name:            "valor hora negativo",
			profissionalID:     1,
			nome:            "Limpeza",
			valorHora:       -10,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "valor_hora",
		},
		{
			name:            "duracao estimada zero",
			profissionalID:     1,
			nome:            "Limpeza",
			valorHora:       50,
			duracaoEstimada: 0,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "duracao_estimada",
		},
		{
			name:            "tipo limpeza invalido",
			profissionalID:     1,
			nome:            "Limpeza",
			valorHora:       50,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpeza("invalido"),
			wantErr:         true,
			wantCampo:       "tipo_limpeza",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l, err := NovaLimpeza(tt.profissionalID, tt.nome, tt.valorHora, tt.duracaoEstimada, tt.tipoLimpeza)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NovaLimpeza() error = nil; want error")
				}
				var ev *ErroValidacao
				if errors.As(err, &ev) && ev.Campo != tt.wantCampo {
					t.Errorf("ErroValidacao.Campo = %q; want %q", ev.Campo, tt.wantCampo)
				}
				return
			}

			if err != nil {
				t.Fatalf("NovaLimpeza() unexpected error: %v", err)
			}
			if l.ProfissionalID != tt.profissionalID {
				t.Errorf("ProfissionalID = %d; want %d", l.ProfissionalID, tt.profissionalID)
			}
			if l.Nome != tt.nome {
				t.Errorf("Nome = %q; want %q", l.Nome, tt.nome)
			}
		})
	}
}

func TestLimpeza_PrecoTotal(t *testing.T) {
	t.Parallel()

	l := &Limpeza{ValorHora: 50, DuracaoEstimada: 3}
	got := l.PrecoTotal()
	want := 150.0

	if got != want {
		t.Errorf("PrecoTotal() = %f; want %f", got, want)
	}
}

func TestLimpeza_EPublicadoPor(t *testing.T) {
	t.Parallel()

	l := &Limpeza{ProfissionalID: 1}

	tests := []struct {
		name        string
		profissionalID int
		want        bool
	}{
		{"id correto", 1, true},
		{"id diferente", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := l.EPublicadoPor(tt.profissionalID)
			if got != tt.want {
				t.Errorf("EPublicadoPor(%d) = %v; want %v", tt.profissionalID, got, tt.want)
			}
		})
	}
}

func TestLimpeza_VerificarPropriedade(t *testing.T) {
	t.Parallel()

	l := &Limpeza{ProfissionalID: 1}

	tests := []struct {
		name        string
		profissionalID int
		wantErr     bool
	}{
		{"dono correto", 1, false},
		{"nao e dono", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := l.VerificarPropriedade(tt.profissionalID)

			if tt.wantErr {
				if !errors.Is(err, errosdominio.ErrNaoEProfissionalDaLimpeza) {
					t.Errorf("VerificarPropriedade() error = %v; want ErrNaoEProfissionalDaLimpeza", err)
				}
				return
			}

			if err != nil {
				t.Errorf("VerificarPropriedade() unexpected error: %v", err)
			}
		})
	}
}
