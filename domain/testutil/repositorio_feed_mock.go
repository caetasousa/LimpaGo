package testutil

import (
	"context"

	"phresh-go/domain/entity"
)

type RepositorioFeedMock struct {
	Itens []*entity.ItemFeed
}

func NovoRepositorioFeedMock() *RepositorioFeedMock {
	return &RepositorioFeedMock{}
}

func (r *RepositorioFeedMock) BuscarPaginaFeed(_ context.Context, pagina, tamanhoPagina int) (*entity.PaginaFeed, error) {
	total := len(r.Itens)
	inicio := (pagina - 1) * tamanhoPagina
	if inicio >= total {
		return &entity.PaginaFeed{
			Itens:         nil,
			TotalItens:    total,
			Pagina:        pagina,
			TamanhoPagina: tamanhoPagina,
		}, nil
	}

	fim := inicio + tamanhoPagina
	if fim > total {
		fim = total
	}

	return &entity.PaginaFeed{
		Itens:         r.Itens[inicio:fim],
		TotalItens:    total,
		Pagina:        pagina,
		TamanhoPagina: tamanhoPagina,
	}, nil
}
