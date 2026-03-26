package repository

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioAvaliacao interface {
	BuscarPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Avaliacao, error)
	ListarPorProfissional(ctx context.Context, profissionalID int) ([]*entity.Avaliacao, error)
	BuscarAgregadoPorProfissional(ctx context.Context, profissionalID int) (*entity.AgregadoAvaliacao, error)
	Salvar(ctx context.Context, avaliacao *entity.Avaliacao) error
}
