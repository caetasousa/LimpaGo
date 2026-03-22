package valueobject

import "testing"

func TestNovaPaginacao(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		pagina            int
		tamanhoPagina     int
		wantPagina        int
		wantTamanhoPagina int
	}{
		{"valores normais", 2, 30, 2, 30},
		{"limites minimos", 1, 1, 1, 1},
		{"limite maximo tamanho", 1, 100, 1, 100},
		{"pagina zero corrige para 1", 0, 20, 1, 20},
		{"pagina negativa corrige para 1", -5, 20, 1, 20},
		{"tamanho zero corrige para 20", 1, 0, 1, 20},
		{"tamanho acima de 100 corrige para 20", 1, 101, 1, 20},
		{"tamanho negativo corrige para 20", 1, -1, 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NovaPaginacao(tt.pagina, tt.tamanhoPagina)

			if got.Pagina != tt.wantPagina {
				t.Errorf("Pagina = %d; want %d", got.Pagina, tt.wantPagina)
			}
			if got.TamanhoPagina != tt.wantTamanhoPagina {
				t.Errorf("TamanhoPagina = %d; want %d", got.TamanhoPagina, tt.wantTamanhoPagina)
			}
		})
	}
}
