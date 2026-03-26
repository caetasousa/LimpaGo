package entity

import "testing"

func TestErroValidacao_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		campo    string
		mensagem string
		want     string
	}{
		{
			name:     "formata mensagem corretamente",
			campo:    "email",
			mensagem: "obrigatório",
			want:     "erro de validação no campo 'email': obrigatório",
		},
		{
			name:     "campo e mensagem diferentes",
			campo:    "nome",
			mensagem: "deve ter pelo menos 3 caracteres",
			want:     "erro de validação no campo 'nome': deve ter pelo menos 3 caracteres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &ErroValidacao{Campo: tt.campo, Mensagem: tt.mensagem}
			got := err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q; want %q", got, tt.want)
			}
		})
	}

	t.Run("implementa interface error", func(t *testing.T) {
		t.Parallel()
		var _ error = &ErroValidacao{}
	})
}
