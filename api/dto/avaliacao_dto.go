package dto

import "limpaGo/domain/entity"

// RequisicaoCriarAvaliacao representa o corpo para criar uma avaliação.
type RequisicaoCriarAvaliacao struct {
	LimpezaID  int    `json:"limpeza_id"`
	Nota       int    `json:"nota"`
	Comentario string `json:"comentario"`
}

// RespostaAvaliacao representa uma avaliação na resposta da API.
type RespostaAvaliacao struct {
	ID          int    `json:"id"`
	LimpezaID   int    `json:"limpeza_id"`
	FaxineiroID int    `json:"faxineiro_id"`
	ClienteID   int    `json:"cliente_id"`
	Nota        int    `json:"nota"`
	Comentario  string `json:"comentario"`
}

// DeAvaliacao converte uma entidade Avaliacao para RespostaAvaliacao.
func DeAvaliacao(a *entity.Avaliacao) RespostaAvaliacao {
	return RespostaAvaliacao{
		ID:          a.ID,
		LimpezaID:   a.LimpezaID,
		FaxineiroID: a.FaxineiroID,
		ClienteID:   a.ClienteID,
		Nota:        int(a.Nota),
		Comentario:  a.Comentario,
	}
}

// DeAvaliacaoLista converte uma lista de Avaliacao para lista de RespostaAvaliacao.
func DeAvaliacaoLista(lista []*entity.Avaliacao) []RespostaAvaliacao {
	resultado := make([]RespostaAvaliacao, len(lista))
	for i, a := range lista {
		resultado[i] = DeAvaliacao(a)
	}
	return resultado
}

// RespostaEstatisticasFaxineiro representa o agregado de avaliações de um faxineiro.
type RespostaEstatisticasFaxineiro struct {
	FaxineiroID     int     `json:"faxineiro_id"`
	MediaNota       float64 `json:"media_nota"`
	TotalAvaliacoes int     `json:"total_avaliacoes"`
}

// DeAgregadoAvaliacao converte um AgregadoAvaliacao para RespostaEstatisticasFaxineiro.
func DeAgregadoAvaliacao(a *entity.AgregadoAvaliacao) RespostaEstatisticasFaxineiro {
	return RespostaEstatisticasFaxineiro{
		FaxineiroID:     a.FaxineiroID,
		MediaNota:       a.MediaNota,
		TotalAvaliacoes: a.TotalAvaliacoes,
	}
}
