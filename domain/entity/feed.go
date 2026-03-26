package entity

import (
	"time"

	"limpaGo/domain/valueobject"
)

// ItemFeed representa um evento de atividade no feed — como serviços publicados ou atualizados por profissionais.
type ItemFeed struct {
	Limpeza     *Limpeza
	TipoEvento  valueobject.TipoEventoFeed
	DataEvento  time.Time
	NumeroLinha int // usado para paginação baseada em cursor
}

// PaginaFeed é um resultado paginado de itens do feed.
type PaginaFeed struct {
	Itens         []*ItemFeed
	TotalItens    int
	Pagina        int
	TamanhoPagina int
}

func (pf *PaginaFeed) TemMais() bool {
	return pf.Pagina*pf.TamanhoPagina < pf.TotalItens
}
