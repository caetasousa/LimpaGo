package entity

import "testing"

func TestNovoPerfil(t *testing.T) {
	t.Parallel()

	p := NovoPerfil(10, "user@test.com", "usuario1")

	if p.UsuarioID != 10 {
		t.Errorf("UsuarioID = %d; want 10", p.UsuarioID)
	}
	if p.Email != "user@test.com" {
		t.Errorf("Email = %q; want %q", p.Email, "user@test.com")
	}
	if p.NomeUsuario != "usuario1" {
		t.Errorf("NomeUsuario = %q; want %q", p.NomeUsuario, "usuario1")
	}
}

func TestNovoPerfilProfissional(t *testing.T) {
	t.Parallel()

	p := NovoPerfilProfissional(5)

	if p.UsuarioID != 5 {
		t.Errorf("UsuarioID = %d; want 5", p.UsuarioID)
	}
	if p.Verificado {
		t.Error("Verificado = true; want false")
	}
}

func TestNovoPerfilCliente(t *testing.T) {
	t.Parallel()

	p := NovoPerfilCliente(7)

	if p.UsuarioID != 7 {
		t.Errorf("UsuarioID = %d; want 7", p.UsuarioID)
	}
}
