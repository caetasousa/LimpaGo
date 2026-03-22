package entity

import (
	"errors"
	"testing"
)

func TestNovoUsuario(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		email       string
		nomeUsuario string
		wantErr     bool
		wantCampo   string
		wantEmail   string
	}{
		{
			name:        "email e username validos",
			email:       "User@Test.COM",
			nomeUsuario: "usuario1",
			wantEmail:   "user@test.com",
		},
		{
			name:        "normaliza email lowercase e trim",
			email:       "  EMAIL@Test.COM  ",
			nomeUsuario: "abc123",
			wantEmail:   "email@test.com",
		},
		{
			name:        "username com tres chars",
			email:       "a@b.com",
			nomeUsuario: "abc",
			wantEmail:   "a@b.com",
		},
		{
			name:        "username com underscore",
			email:       "a@b.com",
			nomeUsuario: "user_name",
			wantEmail:   "a@b.com",
		},
		{
			name:        "username com hifen",
			email:       "a@b.com",
			nomeUsuario: "user-name",
			wantEmail:   "a@b.com",
		},
		{
			name:        "username com maiuscula e numeros",
			email:       "a@b.com",
			nomeUsuario: "User123",
			wantEmail:   "a@b.com",
		},
		{
			name:        "email vazio",
			email:       "",
			nomeUsuario: "user",
			wantErr:     true,
			wantCampo:   "email",
		},
		{
			name:        "email somente espacos",
			email:       "   ",
			nomeUsuario: "user",
			wantErr:     true,
			wantCampo:   "email",
		},
		{
			name:        "username curto",
			email:       "a@b.com",
			nomeUsuario: "ab",
			wantErr:     true,
			wantCampo:   "nome_usuario",
		},
		{
			name:        "username com espaco",
			email:       "a@b.com",
			nomeUsuario: "user name",
			wantErr:     true,
			wantCampo:   "nome_usuario",
		},
		{
			name:        "username com caractere invalido",
			email:       "a@b.com",
			nomeUsuario: "user@name",
			wantErr:     true,
			wantCampo:   "nome_usuario",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u, err := NovoUsuario(tt.email, tt.nomeUsuario)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NovoUsuario() error = nil; want error")
				}
				var ev *ErroValidacao
				if errors.As(err, &ev) && ev.Campo != tt.wantCampo {
					t.Errorf("ErroValidacao.Campo = %q; want %q", ev.Campo, tt.wantCampo)
				}
				return
			}

			if err != nil {
				t.Fatalf("NovoUsuario() unexpected error: %v", err)
			}
			if u.Email != tt.wantEmail {
				t.Errorf("Email = %q; want %q", u.Email, tt.wantEmail)
			}
			if !u.Ativo {
				t.Error("Ativo = false; want true")
			}
		})
	}
}

func TestUsuario_EFaxineiro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		perfilFaxineiro *PerfilFaxineiro
		want            bool
	}{
		{"com perfil faxineiro", &PerfilFaxineiro{}, true},
		{"sem perfil faxineiro", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &Usuario{PerfilFaxineiro: tt.perfilFaxineiro}
			got := u.EFaxineiro()
			if got != tt.want {
				t.Errorf("EFaxineiro() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_ECliente(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		perfilCliente *PerfilCliente
		want          bool
	}{
		{"com perfil cliente", &PerfilCliente{}, true},
		{"sem perfil cliente", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &Usuario{PerfilCliente: tt.perfilCliente}
			got := u.ECliente()
			if got != tt.want {
				t.Errorf("ECliente() = %v; want %v", got, tt.want)
			}
		})
	}
}
