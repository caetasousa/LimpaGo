package valueobject

import "testing"

func TestStatusSolicitacao_Validar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  StatusSolicitacao
		wantErr bool
	}{
		{"pendente valido", StatusSolicitacaoPendente, false},
		{"aceita valido", StatusSolicitacaoAceita, false},
		{"rejeitada valido", StatusSolicitacaoRejeitada, false},
		{"cancelada valido", StatusSolicitacaoCancelada, false},
		{"concluida valido", StatusSolicitacaoConcluida, false},
		{"status invalido", StatusSolicitacao("desconhecido"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.status.Validar()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validar() error = %v; wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStatusSolicitacao_PodeSerCanceladaPeloCliente(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status StatusSolicitacao
		want   bool
	}{
		{"pendente pode cancelar", StatusSolicitacaoPendente, true},
		{"aceita pode cancelar", StatusSolicitacaoAceita, true},
		{"rejeitada nao pode cancelar", StatusSolicitacaoRejeitada, false},
		{"cancelada nao pode cancelar", StatusSolicitacaoCancelada, false},
		{"concluida nao pode cancelar", StatusSolicitacaoConcluida, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.status.PodeSerCanceladaPeloCliente()
			if got != tt.want {
				t.Errorf("PodeSerCanceladaPeloCliente() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestStatusSolicitacao_PodeSerRejeitadaPeloFaxineiro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status StatusSolicitacao
		want   bool
	}{
		{"pendente pode rejeitar", StatusSolicitacaoPendente, true},
		{"aceita nao pode rejeitar", StatusSolicitacaoAceita, false},
		{"rejeitada nao pode rejeitar", StatusSolicitacaoRejeitada, false},
		{"cancelada nao pode rejeitar", StatusSolicitacaoCancelada, false},
		{"concluida nao pode rejeitar", StatusSolicitacaoConcluida, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.status.PodeSerRejeitadaPeloFaxineiro()
			if got != tt.want {
				t.Errorf("PodeSerRejeitadaPeloFaxineiro() = %v; want %v", got, tt.want)
			}
		})
	}
}
