package dto

import (
	"time"

	"limpaGo/domain/entity"
)

// RequisicaoCriarSolicitacao representa o corpo para criar uma solicitação.
type RequisicaoCriarSolicitacao struct {
	LimpezaID    int       `json:"limpeza_id"`
	DataAgendada time.Time `json:"data_agendada"`
}

// RespostaSolicitacao representa uma solicitação na resposta da API.
type RespostaSolicitacao struct {
	ID                int         `json:"id"`
	ClienteID         int         `json:"cliente_id"`
	LimpezaID         int         `json:"limpeza_id"`
	Status            string      `json:"status"`
	DataAgendada      time.Time   `json:"data_agendada"`
	PrecoTotal        float64     `json:"preco_total"`
	MultaCancelamento float64     `json:"multa_cancelamento,omitempty"`
	Endereco          EnderecoDTO `json:"endereco"`
}

// DeSolicitacao converte uma entidade Solicitacao para RespostaSolicitacao.
func DeSolicitacao(s *entity.Solicitacao) RespostaSolicitacao {
	return RespostaSolicitacao{
		ID:                s.ID,
		ClienteID:         s.ClienteID,
		LimpezaID:         s.LimpezaID,
		Status:            string(s.Status),
		DataAgendada:      s.DataAgendada,
		PrecoTotal:        s.PrecoTotal,
		MultaCancelamento: s.MultaCancelamento,
		Endereco:          DeEndereco(s.Endereco),
	}
}

// DeSolicitacaoLista converte uma lista de Solicitacao para lista de RespostaSolicitacao.
func DeSolicitacaoLista(lista []*entity.Solicitacao) []RespostaSolicitacao {
	resultado := make([]RespostaSolicitacao, len(lista))
	for i, s := range lista {
		resultado[i] = DeSolicitacao(s)
	}
	return resultado
}
