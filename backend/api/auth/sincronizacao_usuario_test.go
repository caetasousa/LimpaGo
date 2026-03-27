package auth_test

import (
	"context"
	"errors"
	"testing"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func TestSincronizacaoUsuario_RegistrarCriaUsuario(t *testing.T) {
	t.Parallel()
	repoUsuarios := testutil.NovoRepositorioUsuarioMock()
	repoPerfis := testutil.NovoRepositorioPerfilMock()
	svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)

	ctx := context.Background()

	tests := []struct {
		name        string
		email       string
		nomeUsuario string
		wantErr     bool
		wantErrIs   error
	}{
		{name: "email e nome válidos criam usuário", email: "novo@local.com", nomeUsuario: "novousuario"},
		{name: "email vazio retorna erro de validação", email: "", nomeUsuario: "usuario", wantErr: true},
		{name: "email duplicado retorna ErrEmailJaUtilizado", email: "dup@local.com", nomeUsuario: "dup", wantErr: true, wantErrIs: errosdominio.ErrEmailJaUtilizado},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			localRepo := testutil.NovoRepositorioUsuarioMock()
			localPerfis := testutil.NovoRepositorioPerfilMock()
			localSvc := service.NovoServicoUsuario(localRepo, localPerfis)

			if tt.wantErrIs == errosdominio.ErrEmailJaUtilizado {
				_, _ = localSvc.Registrar(ctx, tt.email, tt.nomeUsuario)
			}
			_ = svcUsuario

			usuario, err := localSvc.Registrar(ctx, tt.email, tt.nomeUsuario)
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
			if usuario.Email != tt.email {
				t.Errorf("got email %q; want %q", usuario.Email, tt.email)
			}
		})
	}
}
