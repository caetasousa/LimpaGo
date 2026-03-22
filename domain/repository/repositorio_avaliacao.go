package repository

import (
	"context"

	"phresh-go/domain/entity"
)

type RepositorioAvaliacao interface {
	BuscarPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Avaliacao, error)
	ListarPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Avaliacao, error)
	BuscarAgregadoPorFaxineiro(ctx context.Context, faxineiroID int) (*entity.AgregadoAvaliacao, error)
	Salvar(ctx context.Context, avaliacao *entity.Avaliacao) error
}
