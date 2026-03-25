package auth_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/api/auth"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func novoServicoSincronizacao(t *testing.T) (*auth.ServicoSincronizacao, *testutil.RepositorioUsuarioMock) {
	t.Helper()
	repoUsuarios := testutil.NovoRepositorioUsuarioMock()
	repoPerfis := testutil.NovoRepositorioPerfilMock()
	svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)
	sinc := auth.NovoServicoSincronizacao(repoUsuarios, svcUsuario)
	return sinc, repoUsuarios
}

func TestSincronizacaoUsuario_SincronizarOuBuscar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		email         string
		nomeUsuario   string
		preRegistrar  bool
		wantErr       bool
		wantErrIs     error
		wantNomeEmail string
	}{
		{
			name:          "primeiro acesso cria usuario automaticamente",
			email:         "novo@zitadel.com",
			nomeUsuario:   "novousuario",
			preRegistrar:  false,
			wantNomeEmail: "novo@zitadel.com",
		},
		{
			name:          "segundo acesso retorna usuario existente sem duplicar",
			email:         "existente@zitadel.com",
			nomeUsuario:   "existente",
			preRegistrar:  true,
			wantNomeEmail: "existente@zitadel.com",
		},
		{
			name:        "email vazio retorna erro de validação",
			email:       "",
			nomeUsuario: "usuario",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sinc, repoUsuarios := novoServicoSincronizacao(t)
			ctx := context.Background()

			if tt.preRegistrar {
				repoPerfis := testutil.NovoRepositorioPerfilMock()
				svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)
				_, err := svcUsuario.Registrar(ctx, tt.email, tt.nomeUsuario)
				if err != nil {
					t.Fatalf("pré-registro falhou: %v", err)
				}
			}

			usuario, err := sinc.SincronizarOuBuscar(ctx, tt.email, tt.nomeUsuario)

			if tt.wantErr {
				if err == nil {
					t.Fatal("got nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("got err %v; want %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if usuario.Email != tt.wantNomeEmail {
				t.Errorf("got email %q; want %q", usuario.Email, tt.wantNomeEmail)
			}
			if usuario.ID == 0 {
				t.Error("got ID 0; want > 0")
			}
		})
	}
}

func TestSincronizacaoUsuario_SincronizarDuasVezesNaoDuplica(t *testing.T) {
	t.Parallel()
	sinc, repoUsuarios := novoServicoSincronizacao(t)
	ctx := context.Background()

	email := "dedup@zitadel.com"
	nome := "dedupuser"

	u1, err := sinc.SincronizarOuBuscar(ctx, email, nome)
	if err != nil {
		t.Fatalf("primeira sincronização falhou: %v", err)
	}

	u2, err := sinc.SincronizarOuBuscar(ctx, email, nome)
	if err != nil {
		t.Fatalf("segunda sincronização falhou: %v", err)
	}

	if u1.ID != u2.ID {
		t.Errorf("segunda sincronização criou usuário duplicado: got IDs %d e %d", u1.ID, u2.ID)
	}

	_ = repoUsuarios // repoUsuarios usado indiretamente via sinc
}

func TestSincronizacaoUsuario_EmailVazioRetornaErroDeValidacao(t *testing.T) {
	t.Parallel()
	sinc, _ := novoServicoSincronizacao(t)
	ctx := context.Background()

	_, err := sinc.SincronizarOuBuscar(ctx, "", "usuario")
	if err == nil {
		t.Fatal("expected error for empty email; got nil")
	}
	// O ServicoUsuario retorna ErrEmailJaUtilizado ou erro de validação
	if errors.Is(err, errosdominio.ErrEmailJaUtilizado) {
		t.Error("unexpected ErrEmailJaUtilizado for empty email")
	}
}
