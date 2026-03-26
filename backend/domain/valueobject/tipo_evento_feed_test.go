package valueobject

import "testing"

func TestTipoEventoFeed_Validar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tipo    TipoEventoFeed
		wantErr bool
	}{
		{"criacao valido", TipoEventoFeedCriacao, false},
		{"atualizacao valido", TipoEventoFeedAtualizacao, false},
		{"tipo invalido", TipoEventoFeed("exclusao"), true},
		{"tipo vazio", TipoEventoFeed(""), true},
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
