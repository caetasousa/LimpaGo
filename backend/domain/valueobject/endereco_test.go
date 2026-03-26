package valueobject

import "testing"

func TestEndereco_EstaPreenchido(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endereco Endereco
		want     bool
	}{
		{
			name:     "rua e cidade preenchidos",
			endereco: Endereco{Rua: "Rua A", Cidade: "São Paulo"},
			want:     true,
		},
		{
			name: "todos campos preenchidos",
			endereco: Endereco{
				Rua: "Rua A", Complemento: "Apto 1", Bairro: "Centro",
				Cidade: "São Paulo", Estado: "SP", CEP: "01000-000",
			},
			want: true,
		},
		{
			name:     "rua vazia",
			endereco: Endereco{Rua: "", Cidade: "São Paulo"},
			want:     false,
		},
		{
			name:     "cidade vazia",
			endereco: Endereco{Rua: "Rua A", Cidade: ""},
			want:     false,
		},
		{
			name:     "struct zerada",
			endereco: Endereco{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.endereco.EstaPreenchido()
			if got != tt.want {
				t.Errorf("EstaPreenchido() = %v; want %v", got, tt.want)
			}
		})
	}
}
