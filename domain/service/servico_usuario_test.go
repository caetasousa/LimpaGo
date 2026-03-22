package service_test

import (
	"context"
	"errors"
	"testing"

	"phresh-go/domain/entity"
	errosdominio "phresh-go/domain/errors"
	"phresh-go/domain/service"
	"phresh-go/domain/testutil"
	"phresh-go/domain/valueobject"
)

func setupServicoUsuario(t *testing.T) (*service.ServicoUsuario, *testutil.RepositorioUsuarioMock, *testutil.RepositorioPerfilMock) {
	t.Helper()
	usuarios := testutil.NovoRepositorioUsuarioMock()
	perfis := testutil.NovoRepositorioPerfilMock()
	svc := service.NovoServicoUsuario(usuarios, perfis)
	return svc, usuarios, perfis
}

func TestServicoUsuario_Registrar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(*service.ServicoUsuario)
		email     string
		username  string
		wantErr   bool
		wantErrIs error
	}{
		{
			name:     "registro valido",
			email:    "user@test.com",
			username: "usuario1",
		},
		{
			name: "email duplicado",
			setup: func(svc *service.ServicoUsuario) {
				_, _ = svc.Registrar(ctx, "dup@test.com", "user1")
			},
			email:     "dup@test.com",
			username:  "user2",
			wantErr:   true,
			wantErrIs: errosdominio.ErrEmailJaUtilizado,
		},
		{
			name: "username duplicado",
			setup: func(svc *service.ServicoUsuario) {
				_, _ = svc.Registrar(ctx, "a@test.com", "mesmousuario")
			},
			email:     "b@test.com",
			username:  "mesmousuario",
			wantErr:   true,
			wantErrIs: errosdominio.ErrNomeUsuarioJaUtilizado,
		},
		{
			name:    "email vazio",
			email:   "",
			username: "user1",
			wantErr: true,
		},
		{
			name:     "username invalido",
			email:    "a@test.com",
			username: "ab",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, _, _ := setupServicoUsuario(t)
			if tt.setup != nil {
				tt.setup(svc)
			}

			u, err := svc.Registrar(ctx, tt.email, tt.username)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Registrar() error = nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("error = %v; want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("Registrar() unexpected error: %v", err)
			}
			if !u.Ativo {
				t.Error("Ativo = false; want true")
			}
			if u.Perfil == nil {
				t.Error("Perfil = nil; want auto-created profile")
			}
			if u.ID == 0 {
				t.Error("ID = 0; want assigned ID")
			}
		})
	}
}

func TestServicoUsuario_CriarPerfilFaxineiro(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("perfil criado", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)

		perfil, err := svc.CriarPerfilFaxineiro(ctx, 1)
		if err != nil {
			t.Fatalf("CriarPerfilFaxineiro() unexpected error: %v", err)
		}
		if perfil.UsuarioID != 1 {
			t.Errorf("UsuarioID = %d; want 1", perfil.UsuarioID)
		}
	})

	t.Run("duplicado", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)
		_, _ = svc.CriarPerfilFaxineiro(ctx, 1)

		_, err := svc.CriarPerfilFaxineiro(ctx, 1)
		if !errors.Is(err, errosdominio.ErrPerfilFaxineiroJaExiste) {
			t.Errorf("error = %v; want ErrPerfilFaxineiroJaExiste", err)
		}
	})
}

func TestServicoUsuario_CriarPerfilCliente(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("perfil criado", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)

		perfil, err := svc.CriarPerfilCliente(ctx, 1)
		if err != nil {
			t.Fatalf("CriarPerfilCliente() unexpected error: %v", err)
		}
		if perfil.UsuarioID != 1 {
			t.Errorf("UsuarioID = %d; want 1", perfil.UsuarioID)
		}
	})

	t.Run("duplicado", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)
		_, _ = svc.CriarPerfilCliente(ctx, 1)

		_, err := svc.CriarPerfilCliente(ctx, 1)
		if !errors.Is(err, errosdominio.ErrPerfilClienteJaExiste) {
			t.Errorf("error = %v; want ErrPerfilClienteJaExiste", err)
		}
	})
}

func TestServicoUsuario_AtualizarPerfil(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc, _, perfis := setupServicoUsuario(t)
	_ = perfis.Salvar(ctx, entity.NovoPerfil(1, "a@b.com", "user1"))

	p, err := svc.AtualizarPerfil(ctx, 1, "João Silva", "11999999999", "http://img.com/foto.jpg")
	if err != nil {
		t.Fatalf("AtualizarPerfil() unexpected error: %v", err)
	}
	if p.NomeCompleto != "João Silva" {
		t.Errorf("NomeCompleto = %q; want %q", p.NomeCompleto, "João Silva")
	}
	if p.Telefone != "11999999999" {
		t.Errorf("Telefone = %q; want %q", p.Telefone, "11999999999")
	}
}

func TestServicoUsuario_AtualizarPerfilFaxineiro(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("atualiza dados profissionais", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)
		_, _ = svc.CriarPerfilFaxineiro(ctx, 1)

		p, err := svc.AtualizarPerfilFaxineiro(ctx, 1, "Profissional experiente", 5, []string{"limpeza_padrao"}, []string{"São Paulo"})
		if err != nil {
			t.Fatalf("AtualizarPerfilFaxineiro() unexpected error: %v", err)
		}
		if p.Descricao != "Profissional experiente" {
			t.Errorf("Descricao = %q; want %q", p.Descricao, "Profissional experiente")
		}
		if p.AnosExperiencia != 5 {
			t.Errorf("AnosExperiencia = %d; want 5", p.AnosExperiencia)
		}
	})

	t.Run("perfil inexistente", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)

		_, err := svc.AtualizarPerfilFaxineiro(ctx, 999, "desc", 1, nil, nil)
		if !errors.Is(err, errosdominio.ErrPerfilFaxineiroNaoEncontrado) {
			t.Errorf("error = %v; want ErrPerfilFaxineiroNaoEncontrado", err)
		}
	})
}

func TestServicoUsuario_AtualizarPerfilCliente(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("atualiza dados cliente", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)
		_, _ = svc.CriarPerfilCliente(ctx, 1)

		end := valueobject.Endereco{Rua: "Rua A", Cidade: "SP"}
		p, err := svc.AtualizarPerfilCliente(ctx, 1, end, valueobject.TipoImovelApartamento, 3, 2, 80, "Tem animais")
		if err != nil {
			t.Fatalf("AtualizarPerfilCliente() unexpected error: %v", err)
		}
		if p.Endereco.Rua != "Rua A" {
			t.Errorf("Endereco.Rua = %q; want %q", p.Endereco.Rua, "Rua A")
		}
		if p.Quartos != 3 {
			t.Errorf("Quartos = %d; want 3", p.Quartos)
		}
	})

	t.Run("perfil inexistente", func(t *testing.T) {
		t.Parallel()
		svc, _, _ := setupServicoUsuario(t)

		_, err := svc.AtualizarPerfilCliente(ctx, 999, valueobject.Endereco{}, "", 0, 0, 0, "")
		if !errors.Is(err, errosdominio.ErrPerfilClienteNaoEncontrado) {
			t.Errorf("error = %v; want ErrPerfilClienteNaoEncontrado", err)
		}
	})
}
