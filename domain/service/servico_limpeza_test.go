package service_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
	"limpaGo/domain/valueobject"
)

func setupServicoLimpeza(t *testing.T) *service.ServicoLimpeza {
	t.Helper()
	return service.NovoServicoLimpeza(testutil.NovoRepositorioLimpezaMock())
}

func TestServicoLimpeza_Criar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		nome        string
		valorHora   float64
		wantErr     bool
	}{
		{"criacao valida", "Limpeza Padrão", 50, false},
		{"nome vazio", "", 50, true},
		{"valor hora zero", "Limpeza", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := setupServicoLimpeza(t)

			l, err := svc.Criar(ctx, 1, tt.nome, "desc", tt.valorHora, 3, valueobject.TipoLimpezaPadrao)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Criar() error = nil; want error")
				}
				var ev *entity.ErroValidacao
				if !errors.As(err, &ev) {
					t.Errorf("error type = %T; want *ErroValidacao", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Criar() unexpected error: %v", err)
			}
			if l.ID == 0 {
				t.Error("ID = 0; want assigned ID")
			}
			if l.Nome != tt.nome {
				t.Errorf("Nome = %q; want %q", l.Nome, tt.nome)
			}
		})
	}
}

func TestServicoLimpeza_Atualizar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("atualiza pelo dono", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		l, _ := svc.Criar(ctx, 1, "Original", "desc", 50, 3, valueobject.TipoLimpezaPadrao)

		atualizado, err := svc.Atualizar(ctx, l.ID, 1, "Novo Nome", "nova desc", 60, 4, valueobject.TipoLimpezaPesada)
		if err != nil {
			t.Fatalf("Atualizar() unexpected error: %v", err)
		}
		if atualizado.Nome != "Novo Nome" {
			t.Errorf("Nome = %q; want %q", atualizado.Nome, "Novo Nome")
		}
		if atualizado.ValorHora != 60 {
			t.Errorf("ValorHora = %f; want 60", atualizado.ValorHora)
		}
	})

	t.Run("campos vazios mantidos", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		l, _ := svc.Criar(ctx, 1, "Original", "desc original", 50, 3, valueobject.TipoLimpezaPadrao)

		atualizado, err := svc.Atualizar(ctx, l.ID, 1, "", "", 0, 0, "")
		if err != nil {
			t.Fatalf("Atualizar() unexpected error: %v", err)
		}
		if atualizado.Nome != "Original" {
			t.Errorf("Nome = %q; want %q (kept)", atualizado.Nome, "Original")
		}
		if atualizado.ValorHora != 50 {
			t.Errorf("ValorHora = %f; want 50 (kept)", atualizado.ValorHora)
		}
	})

	t.Run("atualizar por nao dono", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		l, _ := svc.Criar(ctx, 1, "Limpeza", "desc", 50, 3, valueobject.TipoLimpezaPadrao)

		_, err := svc.Atualizar(ctx, l.ID, 999, "Novo", "", 0, 0, "")
		if !errors.Is(err, errosdominio.ErrNaoEProfissionalDaLimpeza) {
			t.Errorf("error = %v; want ErrNaoEProfissionalDaLimpeza", err)
		}
	})
}

func TestServicoLimpeza_Deletar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("deletar pelo dono", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		l, _ := svc.Criar(ctx, 1, "Limpeza", "desc", 50, 3, valueobject.TipoLimpezaPadrao)

		if err := svc.Deletar(ctx, l.ID, 1); err != nil {
			t.Fatalf("Deletar() unexpected error: %v", err)
		}
	})

	t.Run("deletar por nao dono", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		l, _ := svc.Criar(ctx, 1, "Limpeza", "desc", 50, 3, valueobject.TipoLimpezaPadrao)

		err := svc.Deletar(ctx, l.ID, 999)
		if !errors.Is(err, errosdominio.ErrNaoEProfissionalDaLimpeza) {
			t.Errorf("error = %v; want ErrNaoEProfissionalDaLimpeza", err)
		}
	})
}

func TestServicoLimpeza_ListarCatalogo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("retorna lista paginada", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		_, _ = svc.Criar(ctx, 1, "L1", "", 50, 3, valueobject.TipoLimpezaPadrao)
		_, _ = svc.Criar(ctx, 2, "L2", "", 60, 2, valueobject.TipoLimpezaPesada)

		lista, err := svc.ListarCatalogo(ctx, 1, 10)
		if err != nil {
			t.Fatalf("ListarCatalogo() unexpected error: %v", err)
		}
		if len(lista) != 2 {
			t.Errorf("len(lista) = %d; want 2", len(lista))
		}
	})

	t.Run("paginacao invalida corrigida", func(t *testing.T) {
		t.Parallel()
		svc := setupServicoLimpeza(t)
		_, _ = svc.Criar(ctx, 1, "L1", "", 50, 3, valueobject.TipoLimpezaPadrao)

		lista, err := svc.ListarCatalogo(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListarCatalogo() unexpected error: %v", err)
		}
		if lista == nil {
			t.Error("lista = nil; want non-nil with corrected pagination")
		}
	})
}

func TestServicoLimpeza_ListarPorProfissional(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	svc := setupServicoLimpeza(t)
	_, _ = svc.Criar(ctx, 1, "L1", "", 50, 3, valueobject.TipoLimpezaPadrao)
	_, _ = svc.Criar(ctx, 1, "L2", "", 60, 2, valueobject.TipoLimpezaPesada)
	_, _ = svc.Criar(ctx, 2, "L3", "", 70, 1, valueobject.TipoLimpezaExpress)

	lista, err := svc.ListarPorProfissional(ctx, 1)
	if err != nil {
		t.Fatalf("ListarPorProfissional() unexpected error: %v", err)
	}
	if len(lista) != 2 {
		t.Errorf("len(lista) = %d; want 2", len(lista))
	}
}
