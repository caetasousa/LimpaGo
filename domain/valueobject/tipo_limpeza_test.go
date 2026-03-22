package valueobject

import "testing"

func TestTipoLimpeza_Validar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tipo    TipoLimpeza
		wantErr bool
	}{
		{"padrao valido", TipoLimpezaPadrao, false},
		{"pesada valido", TipoLimpezaPesada, false},
		{"express valido", TipoLimpezaExpress, false},
		{"pre_mudanca valido", TipoLimpezaPreMudanca, false},
		{"pos_obra valido", TipoLimpezaPosObra, false},
		{"comercial valido", TipoLimpezaComercial, false},
		{"passadoria valido", TipoLimpezaPassadoria, false},
		{"tipo invalido", TipoLimpeza("inexistente"), true},
		{"tipo vazio", TipoLimpeza(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.tipo.Validar()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validar() error = %v; wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTipoLimpeza_EResidencial(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tipo TipoLimpeza
		want bool
	}{
		{"padrao e residencial", TipoLimpezaPadrao, true},
		{"pesada e residencial", TipoLimpezaPesada, true},
		{"express e residencial", TipoLimpezaExpress, true},
		{"pre_mudanca e residencial", TipoLimpezaPreMudanca, true},
		{"pos_obra e residencial", TipoLimpezaPosObra, true},
		{"passadoria e residencial", TipoLimpezaPassadoria, true},
		{"comercial nao e residencial", TipoLimpezaComercial, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.tipo.EResidencial()
			if got != tt.want {
				t.Errorf("EResidencial() = %v; want %v", got, tt.want)
			}
		})
	}
}
