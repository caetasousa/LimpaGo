package repository

import (
	"context"

	"phresh-go/domain/entity"
)

type RepositorioFeed interface {
	// BuscarPaginaFeed retorna uma lista paginada de eventos de atividade de limpeza.
	BuscarPaginaFeed(ctx context.Context, pagina, tamanhoPagina int) (*entity.PaginaFeed, error)
}
