package entity

import (
	"errors"
	"testing"

	errosdominio "phresh-go/domain/errors"
	"phresh-go/domain/valueobject"
)

func TestNovaLimpeza(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		faxineiroID     int
		nome            string
		valorHora       float64
		duracaoEstimada float64
		tipoLimpeza     valueobject.TipoLimpeza
		wantErr         bool
		wantCampo       string
	}{
		{
			name:            "parametros validos",
			faxineiroID:     1,
			nome:            "Limpeza Padrão",
			valorHora:       50,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
		},
		{
			name:            "nome vazio",
			faxineiroID:     1,
			nome:            "",
			valorHora:       50,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "nome",
		},
		{
			name:            "valor hora zero",
			faxineiroID:     1,
			nome:            "Limpeza",
			valorHora:       0,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "valor_hora",
		},
		{
			name:            "valor hora negativo",
			faxineiroID:     1,
			nome:            "Limpeza",
			valorHora:       -10,
			duracaoEstimada: 3,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "valor_hora",
		},
		{
			name:            "duracao estimada zero",
			faxineiroID:     1,
			nome:            "Limpeza",
			valorHora:       50,
			duracaoEstimada: 0,
			tipoLimpeza:     valueobject.TipoLimpezaPadrao,
			wantErr:         true,
			wantCampo:       "duracao_estimada",
		},
		{
			name:            "tipo limpeza invalido",
			faxineiroID:     1,
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
			l, err := NovaLimpeza(tt.faxineiroID, tt.nome, tt.valorHora, tt.duracaoEstimada, tt.tipoLimpeza)

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
			if l.FaxineiroID != tt.faxineiroID {
				t.Errorf("FaxineiroID = %d; want %d", l.FaxineiroID, tt.faxineiroID)
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

	l := &Limpeza{FaxineiroID: 1}

	tests := []struct {
		name        string
		faxineiroID int
		want        bool
	}{
		{"id correto", 1, true},
		{"id diferente", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := l.EPublicadoPor(tt.faxineiroID)
			if got != tt.want {
				t.Errorf("EPublicadoPor(%d) = %v; want %v", tt.faxineiroID, got, tt.want)
			}
		})
	}
}

func TestLimpeza_VerificarPropriedade(t *testing.T) {
	t.Parallel()

	l := &Limpeza{FaxineiroID: 1}

	tests := []struct {
		name        string
		faxineiroID int
		wantErr     bool
	}{
		{"dono correto", 1, false},
		{"nao e dono", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := l.VerificarPropriedade(tt.faxineiroID)

			if tt.wantErr {
				if !errors.Is(err, errosdominio.ErrNaoEFaxineiroDaLimpeza) {
					t.Errorf("VerificarPropriedade() error = %v; want ErrNaoEFaxineiroDaLimpeza", err)
				}
				return
			}

			if err != nil {
				t.Errorf("VerificarPropriedade() unexpected error: %v", err)
			}
		})
	}
}
