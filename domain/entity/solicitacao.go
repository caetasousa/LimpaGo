package entity

import (
	"time"

	errosdominio "phresh-go/domain/errors"
	"phresh-go/domain/valueobject"
)

// LimiteMultaCancelamento é o prazo mínimo antes do serviço para cancelar sem multa.
const LimiteMultaCancelamento = 24 * time.Hour

// PercentualMultaCancelamento é o percentual do preço total cobrado como multa.
const PercentualMultaCancelamento = 0.20 // 20%

type Solicitacao struct {
	ID            int
	ClienteID     int
	LimpezaID     int
	Status        valueobject.StatusSolicitacao
	DataAgendada  time.Time // data e hora do início do serviço agendado
	PrecoTotal    float64   // valor/hora × duração estimada, calculado no momento da solicitação
	MultaCancelamento float64 // valor da multa aplicada (0 se não houve multa)
	Endereco          valueobject.Endereco // endereço onde o serviço será realizado
	// Relacionamentos
	Cliente      *Usuario
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

func NovaSolicitacao(clienteID, limpezaID int, limpeza *Limpeza, dataAgendada time.Time) (*Solicitacao, error) {
	// Faxineiro não pode solicitar o próprio serviço
	if limpeza.EPublicadoPor(clienteID) {
		return nil, errosdominio.ErrFaxineiroNaoPodeSolicitarProprio
	}

	if dataAgendada.Before(time.Now()) {
		return nil, errosdominio.ErrAgendamentoNoPassado
	}

	return &Solicitacao{
		ClienteID:    clienteID,
		LimpezaID:    limpezaID,
		Status:       valueobject.StatusSolicitacaoPendente,
		DataAgendada: dataAgendada,
		PrecoTotal:   limpeza.PrecoTotal(),
	}, nil
}

// DataFimEstimada retorna a data/hora estimada de término, baseada na duração do serviço.
func (s *Solicitacao) DataFimEstimada(duracaoHoras float64) time.Time {
	return s.DataAgendada.Add(time.Duration(duracaoHoras * float64(time.Hour)))
}

// DefinirEndereco preenche o endereço onde o serviço será realizado.
func (s *Solicitacao) DefinirEndereco(endereco valueobject.Endereco) {
	s.Endereco = endereco
}

// DefinirEnderecoDoCliente copia o endereço do perfil do cliente para a solicitação.
func (s *Solicitacao) DefinirEnderecoDoCliente(perfil *PerfilCliente) {
	s.Endereco = perfil.Endereco
}

// Aceitar define esta solicitação como aceita pelo faxineiro.
// Apenas solicitações pendentes podem ser aceitas.
func (s *Solicitacao) Aceitar() error {
	if s.Status != valueobject.StatusSolicitacaoPendente {
		return errosdominio.ErrSolicitacaoNaoPodeSerAceita
	}
	s.Status = valueobject.StatusSolicitacaoAceita
	return nil
}

// Rejeitar define esta solicitação como rejeitada pelo faxineiro.
func (s *Solicitacao) Rejeitar() error {
	if !s.Status.PodeSerRejeitadaPeloFaxineiro() {
		return errosdominio.ErrSolicitacaoNaoPodeSerRejeitada
	}
	s.Status = valueobject.StatusSolicitacaoRejeitada
	return nil
}

// Cancelar permite que o cliente cancele sua solicitação (pendente ou aceita).
// Se cancelar com menos de 24h antes do serviço agendado, aplica multa de 20%.
func (s *Solicitacao) Cancelar(agora time.Time) error {
	if !s.Status.PodeSerCanceladaPeloCliente() {
		return errosdominio.ErrSolicitacaoNaoPodeSerCancelada
	}

	// Calcular multa se cancelamento em cima da hora
	if s.Status == valueobject.StatusSolicitacaoAceita && s.DataAgendada.Sub(agora) < LimiteMultaCancelamento {
		s.MultaCancelamento = s.PrecoTotal * PercentualMultaCancelamento
	}

	s.Status = valueobject.StatusSolicitacaoCancelada
	return nil
}

// Concluir marca a solicitação como concluída quando uma avaliação é criada.
func (s *Solicitacao) Concluir() {
	s.Status = valueobject.StatusSolicitacaoConcluida
}
