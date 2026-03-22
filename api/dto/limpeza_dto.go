package dto

import "limpaGo/domain/entity"

// RequisicaoCriarLimpeza representa o corpo para criar um serviço de limpeza.
type RequisicaoCriarLimpeza struct {
	Nome            string  `json:"nome"`
	Descricao       string  `json:"descricao"`
	ValorHora       float64 `json:"valor_hora"`
	DuracaoEstimada float64 `json:"duracao_estimada"`
	TipoLimpeza     string  `json:"tipo_limpeza"`
}

// RequisicaoAtualizarLimpeza representa o corpo para atualizar um serviço de limpeza.
type RequisicaoAtualizarLimpeza struct {
	Nome            string  `json:"nome"`
	Descricao       string  `json:"descricao"`
	ValorHora       float64 `json:"valor_hora"`
	DuracaoEstimada float64 `json:"duracao_estimada"`
	TipoLimpeza     string  `json:"tipo_limpeza"`
}

// RespostaLimpeza representa um serviço de limpeza na resposta da API.
type RespostaLimpeza struct {
	ID              int     `json:"id"`
	Nome            string  `json:"nome"`
	Descricao       string  `json:"descricao"`
	ValorHora       float64 `json:"valor_hora"`
	DuracaoEstimada float64 `json:"duracao_estimada"`
	TipoLimpeza     string  `json:"tipo_limpeza"`
	PrecoTotal      float64 `json:"preco_total"`
	FaxineiroID     int     `json:"faxineiro_id"`
}

// DeLimpeza converte uma entidade Limpeza para RespostaLimpeza.
func DeLimpeza(l *entity.Limpeza) RespostaLimpeza {
	return RespostaLimpeza{
		ID:              l.ID,
		Nome:            l.Nome,
		Descricao:       l.Descricao,
		ValorHora:       l.ValorHora,
		DuracaoEstimada: l.DuracaoEstimada,
		TipoLimpeza:     string(l.TipoLimpeza),
		PrecoTotal:      l.PrecoTotal(),
		FaxineiroID:     l.FaxineiroID,
	}
}

// DeLimpezaLista converte uma lista de Limpeza para lista de RespostaLimpeza.
func DeLimpezaLista(lista []*entity.Limpeza) []RespostaLimpeza {
	resultado := make([]RespostaLimpeza, len(lista))
	for i, l := range lista {
		resultado[i] = DeLimpeza(l)
	}
	return resultado
}
