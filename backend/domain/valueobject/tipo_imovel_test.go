package valueobject

import "testing"

func TestTipoImovel_Validar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tipo    TipoImovel
		wantErr bool
	}{
		{"apartamento valido", TipoImovelApartamento, false},
		{"casa valido", TipoImovelCasa, false},
		{"comercial valido", TipoImovelComercial, false},
		{"tipo invalido", TipoImovel("garagem"), true},
		{"tipo vazio", TipoImovel(""), true},
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
