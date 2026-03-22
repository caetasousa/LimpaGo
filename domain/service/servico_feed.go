package service

import (
	"context"

	"phresh-go/domain/entity"
	"phresh-go/domain/repository"
	"phresh-go/domain/valueobject"
)

type ServicoFeed struct {
	feed repository.RepositorioFeed
}

func NovoServicoFeed(feed repository.RepositorioFeed) *ServicoFeed {
	return &ServicoFeed{feed: feed}
}

// BuscarFeed retorna um feed de atividades paginado de eventos de serviços de limpeza.
func (s *ServicoFeed) BuscarFeed(ctx context.Context, pagina, tamanhoPagina int) (*entity.PaginaFeed, error) {
	p := valueobject.NovaPaginacao(pagina, tamanhoPagina)
	return s.feed.BuscarPaginaFeed(ctx, p.Pagina, p.TamanhoPagina)
}
