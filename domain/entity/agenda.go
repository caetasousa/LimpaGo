package entity

import (
	"time"

	errosdominio "limpaGo/domain/errors"
)

// Disponibilidade representa um bloco de horário disponível na agenda do faxineiro.
type Disponibilidade struct {
	ID           int
	FaxineiroID  int
	DiaSemana    time.Weekday // 0=Domingo, 1=Segunda, ..., 6=Sábado
	HoraInicio   int          // hora de início (0-23)
	HoraFim      int          // hora de fim (1-24)
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

func NovaDisponibilidade(faxineiroID int, diaSemana time.Weekday, horaInicio, horaFim int) (*Disponibilidade, error) {
	if horaInicio < 0 || horaInicio > 23 {
		return nil, &ErroValidacao{Campo: "hora_inicio", Mensagem: "hora de início deve estar entre 0 e 23"}
	}
	if horaFim < 1 || horaFim > 24 {
		return nil, &ErroValidacao{Campo: "hora_fim", Mensagem: "hora de fim deve estar entre 1 e 24"}
	}
	if horaFim <= horaInicio {
		return nil, &ErroValidacao{Campo: "hora_fim", Mensagem: "hora de fim deve ser maior que hora de início"}
	}

	return &Disponibilidade{
		FaxineiroID: faxineiroID,
		DiaSemana:   diaSemana,
		HoraInicio:  horaInicio,
		HoraFim:     horaFim,
	}, nil
}

// DuracaoHoras retorna quantas horas este bloco de disponibilidade tem.
func (d *Disponibilidade) DuracaoHoras() int {
	return d.HoraFim - d.HoraInicio
}

// Bloqueio representa um horário ocupado na agenda do faxineiro.
// Pode ser gerado por uma solicitação aceita ou pelo próprio faxineiro.
type Bloqueio struct {
	ID            int
	FaxineiroID   int
	SolicitacaoID *int // nil = bloqueio pessoal do faxineiro, preenchido = serviço agendado
	DataInicio    time.Time
	DataFim       time.Time
	CriadoEm      time.Time
}

// EPessoal retorna true se o bloqueio é pessoal (não vinculado a uma solicitação).
func (b *Bloqueio) EPessoal() bool {
	return b.SolicitacaoID == nil
}

// NovoBloqueioServico cria um bloqueio gerado pela aceitação de uma solicitação.
func NovoBloqueioServico(faxineiroID, solicitacaoID int, dataInicio, dataFim time.Time) (*Bloqueio, error) {
	if err := validarPeriodoBloqueio(dataInicio, dataFim); err != nil {
		return nil, err
	}

	return &Bloqueio{
		FaxineiroID:   faxineiroID,
		SolicitacaoID: &solicitacaoID,
		DataInicio:    dataInicio,
		DataFim:       dataFim,
	}, nil
}

// NovoBloqueiopessoal cria um bloqueio pessoal do faxineiro (ex: consulta médica, dentista, folga).
func NovoBloqueiopessoal(faxineiroID int, dataInicio, dataFim time.Time) (*Bloqueio, error) {
	if err := validarPeriodoBloqueio(dataInicio, dataFim); err != nil {
		return nil, err
	}

	return &Bloqueio{
		FaxineiroID: faxineiroID,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
	}, nil
}

func validarPeriodoBloqueio(dataInicio, dataFim time.Time) error {
	if dataFim.Before(dataInicio) || dataFim.Equal(dataInicio) {
		return &ErroValidacao{Campo: "data_fim", Mensagem: "data de fim deve ser posterior à data de início"}
	}
	if dataInicio.Before(time.Now()) {
		return errosdominio.ErrAgendamentoNoPassado
	}
	return nil
}
