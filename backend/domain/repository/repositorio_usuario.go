package repository

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioUsuario interface {
	BuscarPorID(ctx context.Context, id int) (*entity.Usuario, error)
	BuscarPorEmail(ctx context.Context, email string) (*entity.Usuario, error)
	BuscarPorNomeUsuario(ctx context.Context, nomeUsuario string) (*entity.Usuario, error)
	Salvar(ctx context.Context, usuario *entity.Usuario) error
}
