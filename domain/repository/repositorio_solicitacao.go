package repository

import (
	"context"

	"phresh-go/domain/entity"
)

type RepositorioSolicitacao interface {
	BuscarPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error)
	BuscarAtivaPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) // pendente ou aceita
	ListarPorLimpeza(ctx context.Context, limpezaID int) ([]*entity.Solicitacao, error)
	ListarPorCliente(ctx context.Context, clienteID int) ([]*entity.Solicitacao, error)
	Salvar(ctx context.Context, solicitacao *entity.Solicitacao) error
	Atualizar(ctx context.Context, solicitacao *entity.Solicitacao) error
	Deletar(ctx context.Context, clienteID, limpezaID int) error
}
