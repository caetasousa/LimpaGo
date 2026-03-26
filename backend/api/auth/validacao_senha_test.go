package auth_test

import (
	"errors"
	"testing"

	"limpaGo/api/auth"
)

func TestValidarForcaSenha(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		senha     string
		wantErr   bool
		wantErrIs error
	}{
		{name: "válida", senha: "Senha123", wantErr: false},
		{name: "válida longa", senha: "MinhaS3nhaSegura!", wantErr: false},
		{name: "muito curta", senha: "Ab1", wantErr: true, wantErrIs: auth.ErrSenhaFraca},
		{name: "sem maiúscula", senha: "senha123", wantErr: true, wantErrIs: auth.ErrSenhaFraca},
		{name: "sem dígito", senha: "SenhaForte", wantErr: true, wantErrIs: auth.ErrSenhaFraca},
		{name: "apenas letras minúsculas", senha: "senhasemtudo", wantErr: true, wantErrIs: auth.ErrSenhaFraca},
		{name: "vazia", senha: "", wantErr: true, wantErrIs: auth.ErrSenhaFraca},
		{name: "exatamente 8 chars válida", senha: "Senha12!", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := auth.ValidarForcaSenha(tt.senha)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error; got nil")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("got err %v; want %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
