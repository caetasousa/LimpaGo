package valueobject

import "errors"

type StatusSolicitacao string

const (
	StatusSolicitacaoPendente  StatusSolicitacao = "pendente"
	StatusSolicitacaoAceita    StatusSolicitacao = "aceita"
	StatusSolicitacaoRejeitada StatusSolicitacao = "rejeitada"
	StatusSolicitacaoCancelada StatusSolicitacao = "cancelada"
	StatusSolicitacaoConcluida StatusSolicitacao = "concluida"
)

func (s StatusSolicitacao) Validar() error {
	switch s {
	case StatusSolicitacaoPendente, StatusSolicitacaoAceita, StatusSolicitacaoRejeitada, StatusSolicitacaoCancelada, StatusSolicitacaoConcluida:
		return nil
	default:
		return errors.New("status de solicitação inválido")
	}
}

// PodeSerCanceladaPeloCliente retorna true se o cliente pode cancelar a solicitação.
// O cliente pode cancelar solicitações pendentes ou aceitas.
func (s StatusSolicitacao) PodeSerCanceladaPeloCliente() bool {
	return s == StatusSolicitacaoPendente || s == StatusSolicitacaoAceita
}

// PodeSerRejeitadaPeloFaxineiro retorna true se o faxineiro pode rejeitar a solicitação.
// Apenas solicitações pendentes podem ser rejeitadas.
func (s StatusSolicitacao) PodeSerRejeitadaPeloFaxineiro() bool {
	return s == StatusSolicitacaoPendente
}
