package repository

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioLimpeza interface {
	BuscarPorID(ctx context.Context, id int) (*entity.Limpeza, error)
	ListarPorProfissional(ctx context.Context, profissionalID int) ([]*entity.Limpeza, error)
	ListarTodas(ctx context.Context, pagina, tamanhoPagina int) ([]*entity.Limpeza, error)
	Salvar(ctx context.Context, limpeza *entity.Limpeza) error
	Atualizar(ctx context.Context, limpeza *entity.Limpeza) error
	Deletar(ctx context.Context, id int) error
}
