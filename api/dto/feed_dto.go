package dto

import (
	"time"

	"limpaGo/domain/entity"
)

// RespostaItemFeed representa um item do feed na resposta da API.
type RespostaItemFeed struct {
	Limpeza     RespostaLimpeza `json:"limpeza"`
	TipoEvento  string          `json:"tipo_evento"`
	DataEvento  time.Time       `json:"data_evento"`
	NumeroLinha int             `json:"numero_linha"`
}

// DeItemFeed converte um entity.ItemFeed para RespostaItemFeed.
func DeItemFeed(item *entity.ItemFeed) RespostaItemFeed {
	return RespostaItemFeed{
		Limpeza:     DeLimpeza(item.Limpeza),
		TipoEvento:  string(item.TipoEvento),
		DataEvento:  item.DataEvento,
		NumeroLinha: item.NumeroLinha,
	}
}

// RespostaPaginaFeed representa uma página do feed na resposta da API.
type RespostaPaginaFeed struct {
	Itens        []RespostaItemFeed `json:"itens"`
	TotalItens   int                `json:"total_itens"`
	Pagina       int                `json:"pagina"`
	TamanhoPagina int               `json:"tamanho_pagina"`
	TemMais      bool               `json:"tem_mais"`
}

// DePaginaFeed converte um entity.PaginaFeed para RespostaPaginaFeed.
func DePaginaFeed(p *entity.PaginaFeed) RespostaPaginaFeed {
	itens := make([]RespostaItemFeed, len(p.Itens))
	for i, item := range p.Itens {
		itens[i] = DeItemFeed(item)
	}
	return RespostaPaginaFeed{
		Itens:         itens,
		TotalItens:    p.TotalItens,
		Pagina:        p.Pagina,
		TamanhoPagina: p.TamanhoPagina,
		TemMais:       p.TemMais(),
	}
}
