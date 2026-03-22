package entity

import (
	"errors"
	"testing"
	"time"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

func limpezaParaTeste(faxineiroID int) *Limpeza {
	return &Limpeza{
		ID:              1,
		FaxineiroID:     faxineiroID,
		Nome:            "Limpeza Teste",
		ValorHora:       50,
		DuracaoEstimada: 3,
		TipoLimpeza:     valueobject.TipoLimpezaPadrao,
	}
}

func TestNovaSolicitacao(t *testing.T) {
	t.Parallel()

	limpeza := limpezaParaTeste(1)
	futuro := time.Now().Add(48 * time.Hour)
	passado := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name      string
		clienteID int
		data      time.Time
		wantErr   bool
		wantErrIs error
	}{
		{"cliente valido com data futura", 2, futuro, false, nil},
		{"faxineiro solicitando proprio servico", 1, futuro, true, errosdominio.ErrFaxineiroNaoPodeSolicitarProprio},
		{"data no passado", 2, passado, true, errosdominio.ErrAgendamentoNoPassado},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NovaSolicitacao(tt.clienteID, 1, limpeza, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatal("NovaSolicitacao() error = nil; want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("error = %v; want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("NovaSolicitacao() unexpected error: %v", err)
			}
			if s.Status != valueobject.StatusSolicitacaoPendente {
				t.Errorf("Status = %q; want pendente", s.Status)
			}
			if s.PrecoTotal != 150 {
				t.Errorf("PrecoTotal = %f; want 150", s.PrecoTotal)
			}
		})
	}
}

func TestSolicitacao_DataFimEstimada(t *testing.T) {
	t.Parallel()

	agora := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	s := &Solicitacao{DataAgendada: agora}
	got := s.DataFimEstimada(3)
	want := agora.Add(3 * time.Hour)

	if !got.Equal(want) {
		t.Errorf("DataFimEstimada(3) = %v; want %v", got, want)
	}
}

func TestSolicitacao_DefinirEndereco(t *testing.T) {
	t.Parallel()

	s := &Solicitacao{}
	end := valueobject.Endereco{Rua: "Rua A", Cidade: "SP"}
	s.DefinirEndereco(end)

	if s.Endereco.Rua != "Rua A" || s.Endereco.Cidade != "SP" {
		t.Errorf("Endereco = %+v; want Rua=Rua A, Cidade=SP", s.Endereco)
	}
}

func TestSolicitacao_DefinirEnderecoDoCliente(t *testing.T) {
	t.Parallel()

	s := &Solicitacao{}
	perfil := &PerfilCliente{
		Endereco: valueobject.Endereco{Rua: "Rua B", Cidade: "RJ"},
	}
	s.DefinirEnderecoDoCliente(perfil)

	if s.Endereco.Rua != "Rua B" || s.Endereco.Cidade != "RJ" {
		t.Errorf("Endereco = %+v; want Rua=Rua B, Cidade=RJ", s.Endereco)
	}
}

func TestSolicitacao_Aceitar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  valueobject.StatusSolicitacao
		wantErr bool
	}{
		{"pendente para aceita", valueobject.StatusSolicitacaoPendente, false},
		{"aceita nao pode aceitar", valueobject.StatusSolicitacaoAceita, true},
		{"rejeitada nao pode aceitar", valueobject.StatusSolicitacaoRejeitada, true},
		{"cancelada nao pode aceitar", valueobject.StatusSolicitacaoCancelada, true},
		{"concluida nao pode aceitar", valueobject.StatusSolicitacaoConcluida, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Solicitacao{Status: tt.status}
			err := s.Aceitar()

			if tt.wantErr {
				if !errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerAceita) {
					t.Errorf("Aceitar() error = %v; want ErrSolicitacaoNaoPodeSerAceita", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Aceitar() unexpected error: %v", err)
			}
			if s.Status != valueobject.StatusSolicitacaoAceita {
				t.Errorf("Status = %q; want aceita", s.Status)
			}
		})
	}
}

func TestSolicitacao_Rejeitar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  valueobject.StatusSolicitacao
		wantErr bool
	}{
		{"pendente para rejeitada", valueobject.StatusSolicitacaoPendente, false},
		{"aceita nao pode rejeitar", valueobject.StatusSolicitacaoAceita, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Solicitacao{Status: tt.status}
			err := s.Rejeitar()

			if tt.wantErr {
				if !errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerRejeitada) {
					t.Errorf("Rejeitar() error = %v; want ErrSolicitacaoNaoPodeSerRejeitada", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Rejeitar() unexpected error: %v", err)
			}
			if s.Status != valueobject.StatusSolicitacaoRejeitada {
				t.Errorf("Status = %q; want rejeitada", s.Status)
			}
		})
	}
}

func TestSolicitacao_Cancelar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		status    valueobject.StatusSolicitacao
		preco     float64
		dataAg    time.Time
		agora     time.Time
		wantErr   bool
		wantMulta float64
	}{
		{
			name:    "pendente sem multa",
			status:  valueobject.StatusSolicitacaoPendente,
			preco:   150,
			dataAg:  time.Now().Add(48 * time.Hour),
			agora:   time.Now(),
			wantErr: false,
		},
		{
			name:    "aceita mais de 24h sem multa",
			status:  valueobject.StatusSolicitacaoAceita,
			preco:   150,
			dataAg:  time.Now().Add(48 * time.Hour),
			agora:   time.Now(),
			wantErr: false,
		},
		{
			name:      "aceita menos de 24h com multa 20%",
			status:    valueobject.StatusSolicitacaoAceita,
			preco:     150,
			dataAg:    time.Now().Add(12 * time.Hour),
			agora:     time.Now(),
			wantErr:   false,
			wantMulta: 30,
		},
		{
			name:    "rejeitada nao pode cancelar",
			status:  valueobject.StatusSolicitacaoRejeitada,
			wantErr: true,
		},
		{
			name:    "concluida nao pode cancelar",
			status:  valueobject.StatusSolicitacaoConcluida,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Solicitacao{
				Status:       tt.status,
				PrecoTotal:   tt.preco,
				DataAgendada: tt.dataAg,
			}
			err := s.Cancelar(tt.agora)

			if tt.wantErr {
				if !errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerCancelada) {
					t.Errorf("Cancelar() error = %v; want ErrSolicitacaoNaoPodeSerCancelada", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Cancelar() unexpected error: %v", err)
			}
			if s.Status != valueobject.StatusSolicitacaoCancelada {
				t.Errorf("Status = %q; want cancelada", s.Status)
			}
			if s.MultaCancelamento != tt.wantMulta {
				t.Errorf("MultaCancelamento = %f; want %f", s.MultaCancelamento, tt.wantMulta)
			}
		})
	}
}

func TestSolicitacao_Concluir(t *testing.T) {
	t.Parallel()

	s := &Solicitacao{Status: valueobject.StatusSolicitacaoAceita}
	s.Concluir()

	if s.Status != valueobject.StatusSolicitacaoConcluida {
		t.Errorf("Status = %q; want concluida", s.Status)
	}
}
