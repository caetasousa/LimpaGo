package dto

import (
	"time"

	"limpaGo/domain/entity"
)

// RequisicaoDisponibilidade representa o corpo para adicionar disponibilidade.
type RequisicaoDisponibilidade struct {
	DiaSemana  int `json:"dia_semana"`  // 0=Domingo ... 6=Sábado
	HoraInicio int `json:"hora_inicio"` // 0-23
	HoraFim    int `json:"hora_fim"`    // 1-24
}

// RespostaDisponibilidade representa uma disponibilidade na resposta da API.
type RespostaDisponibilidade struct {
	ID          int `json:"id"`
	ProfissionalID int `json:"profissional_id"`
	DiaSemana   int `json:"dia_semana"`
	HoraInicio  int `json:"hora_inicio"`
	HoraFim     int `json:"hora_fim"`
}

// DeDisponibilidade converte uma entidade Disponibilidade para RespostaDisponibilidade.
func DeDisponibilidade(d *entity.Disponibilidade) RespostaDisponibilidade {
	return RespostaDisponibilidade{
		ID:          d.ID,
		ProfissionalID: d.ProfissionalID,
		DiaSemana:   int(d.DiaSemana),
		HoraInicio:  d.HoraInicio,
		HoraFim:     d.HoraFim,
	}
}

// DeDisponibilidadeLista converte uma lista de Disponibilidade para lista de RespostaDisponibilidade.
func DeDisponibilidadeLista(lista []*entity.Disponibilidade) []RespostaDisponibilidade {
	resultado := make([]RespostaDisponibilidade, len(lista))
	for i, d := range lista {
		resultado[i] = DeDisponibilidade(d)
	}
	return resultado
}

// RequisicaoBloqueio representa o corpo para criar um bloqueio pessoal.
type RequisicaoBloqueio struct {
	DataInicio time.Time `json:"data_inicio"`
	DataFim    time.Time `json:"data_fim"`
}

// RespostaBloqueio representa um bloqueio na resposta da API.
type RespostaBloqueio struct {
	ID            int        `json:"id"`
	ProfissionalID   int        `json:"profissional_id"`
	SolicitacaoID *int       `json:"solicitacao_id,omitempty"`
	DataInicio    time.Time  `json:"data_inicio"`
	DataFim       time.Time  `json:"data_fim"`
	EPessoal      bool       `json:"e_pessoal"`
}

// DeBloqueio converte uma entidade Bloqueio para RespostaBloqueio.
func DeBloqueio(b *entity.Bloqueio) RespostaBloqueio {
	return RespostaBloqueio{
		ID:            b.ID,
		ProfissionalID:   b.ProfissionalID,
		SolicitacaoID: b.SolicitacaoID,
		DataInicio:    b.DataInicio,
		DataFim:       b.DataFim,
		EPessoal:      b.EPessoal(),
	}
}

// DeBloqueioLista converte uma lista de Bloqueio para lista de RespostaBloqueio.
func DeBloqueioLista(lista []*entity.Bloqueio) []RespostaBloqueio {
	resultado := make([]RespostaBloqueio, len(lista))
	for i, b := range lista {
		resultado[i] = DeBloqueio(b)
	}
	return resultado
}
