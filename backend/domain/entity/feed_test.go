package entity

import "testing"

func TestPaginaFeed_TemMais(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		pagina        int
		tamanhoPagina int
		totalItens    int
		want          bool
	}{
		{"tem mais paginas", 1, 10, 25, true},
		{"exatamente uma pagina", 1, 10, 10, false},
		{"ultima pagina", 3, 10, 25, false},
		{"total zero", 1, 10, 0, false},
		{"pagina 2 de 3", 2, 10, 25, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pf := &PaginaFeed{
				Pagina:        tt.pagina,
				TamanhoPagina: tt.tamanhoPagina,
				TotalItens:    tt.totalItens,
			}
			got := pf.TemMais()
			if got != tt.want {
				t.Errorf("TemMais() = %v; want %v", got, tt.want)
			}
		})
	}
}
