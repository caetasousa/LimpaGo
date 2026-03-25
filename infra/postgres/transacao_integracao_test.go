//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	"limpaGo/infra/postgres"
)

func TestComTransacao_Commit(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })

	ctx := context.Background()
	repoUsuario := postgres.NovoRepositorioUsuarioPG(db)
	repoPerfil := postgres.NovoRepositorioPerfilPG(db)

	err := postgres.ComTransacao(ctx, db, func(ctx context.Context) error {
		u := &entity.Usuario{Email: "tx@commit.com", NomeUsuario: "txcommit", Ativo: true}
		if err := repoUsuario.Salvar(ctx, u); err != nil {
			return err
		}
		p := entity.NovoPerfil(u.ID, u.Email, u.NomeUsuario)
		return repoPerfil.Salvar(ctx, p)
	})

	if err != nil {
		t.Fatalf("ComTransacao() commit error: %v", err)
	}

	// Verificar que o usuário e o perfil foram persistidos
	got, err := repoUsuario.BuscarPorEmail(ctx, "tx@commit.com")
	if err != nil {
		t.Fatalf("BuscarPorEmail() error: %v", err)
	}
	if got == nil {
		t.Fatal("usuario nao persistido apos commit")
	}

	perfil, err := repoPerfil.BuscarPorUsuarioID(ctx, got.ID)
	if err != nil {
		t.Fatalf("BuscarPorUsuarioID() error: %v", err)
	}
	if perfil == nil {
		t.Fatal("perfil nao persistido apos commit")
	}
}

func TestComTransacao_Rollback(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })

	ctx := context.Background()
	repoUsuario := postgres.NovoRepositorioUsuarioPG(db)
	erroForcado := errors.New("erro forçado para rollback")

	err := postgres.ComTransacao(ctx, db, func(ctx context.Context) error {
		u := &entity.Usuario{Email: "tx@rollback.com", NomeUsuario: "txrollback", Ativo: true}
		if err := repoUsuario.Salvar(ctx, u); err != nil {
			return err
		}
		return erroForcado
	})

	if !errors.Is(err, erroForcado) {
		t.Errorf("got %v; want %v", err, erroForcado)
	}

	// Verificar que o rollback desfez a inserção
	got, err := repoUsuario.BuscarPorEmail(ctx, "tx@rollback.com")
	if err != nil {
		t.Fatalf("BuscarPorEmail() error: %v", err)
	}
	if got != nil {
		t.Error("usuario persistido apos rollback; want nil")
	}
}

func TestComTransacao_PropagacaoContexto(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })

	ctx := context.Background()
	repoUsuario := postgres.NovoRepositorioUsuarioPG(db)

	// Dois usuários com emails diferentes dentro da mesma transação
	err := postgres.ComTransacao(ctx, db, func(ctx context.Context) error {
		u1 := &entity.Usuario{Email: "tx1@propagacao.com", NomeUsuario: "txprop1", Ativo: true}
		if err := repoUsuario.Salvar(ctx, u1); err != nil {
			return err
		}
		u2 := &entity.Usuario{Email: "tx2@propagacao.com", NomeUsuario: "txprop2", Ativo: true}
		return repoUsuario.Salvar(ctx, u2)
	})

	if err != nil {
		t.Fatalf("ComTransacao() error: %v", err)
	}

	u1, _ := repoUsuario.BuscarPorEmail(ctx, "tx1@propagacao.com")
	u2, _ := repoUsuario.BuscarPorEmail(ctx, "tx2@propagacao.com")
	if u1 == nil || u2 == nil {
		t.Error("ambos usuarios devem estar persistidos apos commit")
	}
}
