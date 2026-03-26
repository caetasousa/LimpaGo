//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
	"limpaGo/infra/postgres"
)

func TestServico_ProfissionalPublicaNovoServicoDeClimpeza(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax@limpeza.com", "faxlimpeza")

	l := &entity.Limpeza{
		ProfissionalID:     profissionalID,
		Nome:            "Limpeza Residencial",
		Descricao:       "Limpeza completa",
		ValorHora:       50.0,
		DuracaoEstimada: 3.0,
		TipoLimpeza:     valueobject.TipoLimpezaPadrao,
	}

	if err := repo.Salvar(ctx, l); err != nil {
		t.Fatalf("Salvar() error: %v", err)
	}
	if l.ID == 0 {
		t.Error("ID = 0; want > 0")
	}
	if l.CriadoEm.IsZero() {
		t.Error("CriadoEm zerado; want preenchido")
	}
}

func TestServico_ConsultarDetalhesDeUmServicoPeloID(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax2@limpeza.com", "faxlimpeza2")
	id := inserirLimpeza(t, db, profissionalID, "Limpeza Busca")

	tests := []struct {
		name    string
		id      int
		wantErr error
		wantNil bool
	}{
		{name: "servico existente retorna todos os detalhes", id: id},
		{
			name:    "servico inexistente retorna erro de nao encontrado",
			id:      999999,
			wantErr: errosdominio.ErrLimpezaNaoEncontrada,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.BuscarPorID(ctx, tt.id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got %v; want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("BuscarPorID() error: %v", err)
			}
			if got == nil {
				t.Fatal("got nil; want limpeza")
			}
			if got.Nome != "Limpeza Busca" {
				t.Errorf("Nome = %q; want %q", got.Nome, "Limpeza Busca")
			}
		})
	}
}

func TestServico_ListarTodosServicosDeUmProfissional(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax3@limpeza.com", "faxlimpeza3")
	inserirLimpeza(t, db, profissionalID, "Servico 1")
	inserirLimpeza(t, db, profissionalID, "Servico 2")

	got, err := repo.ListarPorProfissional(ctx, profissionalID)
	if err != nil {
		t.Fatalf("ListarPorProfissional() error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d; want 2", len(got))
	}
}

func TestServico_ListarTodosServicosComPaginacao(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax4@limpeza.com", "faxlimpeza4")
	for i := 0; i < 5; i++ {
		inserirLimpeza(t, db, profissionalID, "Servico")
	}

	tests := []struct {
		name          string
		pagina        int
		tamanhoPagina int
		wantLen       int
	}{
		{name: "primeira pagina retorna 3 de 5 servicos", pagina: 1, tamanhoPagina: 3, wantLen: 3},
		{name: "segunda pagina retorna os 2 servicos restantes", pagina: 2, tamanhoPagina: 3, wantLen: 2},
		{name: "terceira pagina retorna vazio quando nao ha mais servicos", pagina: 3, tamanhoPagina: 3, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.ListarTodas(ctx, tt.pagina, tt.tamanhoPagina)
			if err != nil {
				t.Fatalf("ListarTodas() error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("len = %d; want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestServico_ProfissionalAtualizaValorEDescricaoDoServico(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax5@limpeza.com", "faxlimpeza5")
	id := inserirLimpeza(t, db, profissionalID, "Original")

	l := &entity.Limpeza{
		ID:              id,
		Nome:            "Atualizado",
		Descricao:       "Descrição atualizada",
		ValorHora:       75.0,
		DuracaoEstimada: 4.0,
		TipoLimpeza:     valueobject.TipoLimpezaPesada,
	}
	if err := repo.Atualizar(ctx, l); err != nil {
		t.Fatalf("Atualizar() error: %v", err)
	}
	if l.AtualizadoEm.IsZero() {
		t.Error("AtualizadoEm zerado; want preenchido")
	}

	got, _ := repo.BuscarPorID(ctx, id)
	if got.Nome != "Atualizado" {
		t.Errorf("Nome = %q; want %q", got.Nome, "Atualizado")
	}
}

func TestServico_ProfissionalRemoveServicoDaPlataforma(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioLimpezaPG(db)
	ctx := context.Background()

	profissionalID := inserirUsuario(t, db, "fax6@limpeza.com", "faxlimpeza6")
	id := inserirLimpeza(t, db, profissionalID, "Para deletar")

	if err := repo.Deletar(ctx, id); err != nil {
		t.Fatalf("Deletar() error: %v", err)
	}

	_, err := repo.BuscarPorID(ctx, id)
	if !errors.Is(err, errosdominio.ErrLimpezaNaoEncontrada) {
		t.Errorf("apos deletar: got %v; want %v", err, errosdominio.ErrLimpezaNaoEncontrada)
	}
}
