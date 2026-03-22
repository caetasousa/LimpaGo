package service

import (
	"context"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/repository"
	"limpaGo/domain/valueobject"
)

// ServicoSolicitacao contém as regras do ciclo de vida das solicitações.
//
// Ciclo de vida da solicitação:
//
//	pendente  → aceita     (faxineiro aceita → bloqueia agenda)
//	pendente  → rejeitada  (faxineiro rejeita)
//	pendente  → cancelada  (cliente cancela)
//	aceita    → concluída  (cliente cria avaliação)
//	aceita    → cancelada  (cliente cancela → libera agenda)
type ServicoSolicitacao struct {
	solicitacoes repository.RepositorioSolicitacao
	limpezas     repository.RepositorioLimpeza
	agenda       *ServicoAgenda
}

func NovoServicoSolicitacao(
	solicitacoes repository.RepositorioSolicitacao,
	limpezas repository.RepositorioLimpeza,
	agenda *ServicoAgenda,
) *ServicoSolicitacao {
	return &ServicoSolicitacao{
		solicitacoes: solicitacoes,
		limpezas:     limpezas,
		agenda:       agenda,
	}
}

// CriarSolicitacao permite que um cliente solicite o serviço de um faxineiro em uma data/hora específica.
// Verifica se o faxineiro está disponível no horário solicitado.
func (s *ServicoSolicitacao) CriarSolicitacao(ctx context.Context, clienteID, limpezaID int, dataAgendada time.Time) (*entity.Solicitacao, error) {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}

	// Verificação de solicitação duplicada (apenas pendentes ou aceitas)
	if existente, err := s.solicitacoes.BuscarAtivaPorClienteELimpeza(ctx, clienteID, limpezaID); err != nil {
		return nil, err
	} else if existente != nil {
		return nil, errosdominio.ErrSolicitacaoDuplicada
	}

	// Calcular fim estimado baseado na duração do serviço
	dataFim := dataAgendada.Add(time.Duration(limpeza.DuracaoEstimada * float64(time.Hour)))

	// Verificar se o faxineiro está disponível neste horário
	if err := s.agenda.VerificarDisponibilidade(ctx, limpeza.FaxineiroID, dataAgendada, dataFim); err != nil {
		return nil, err
	}

	solicitacao, err := entity.NovaSolicitacao(clienteID, limpezaID, limpeza, dataAgendada)
	if err != nil {
		return nil, err
	}

	if err := s.solicitacoes.Salvar(ctx, solicitacao); err != nil {
		return nil, err
	}
	return solicitacao, nil
}

// AceitarSolicitacao permite que o faxineiro aceite a solicitação de um cliente.
// Ao aceitar, o horário é bloqueado na agenda do faxineiro.
func (s *ServicoSolicitacao) AceitarSolicitacao(ctx context.Context, faxineiroID, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}
	if err := limpeza.VerificarPropriedade(faxineiroID); err != nil {
		return nil, errosdominio.ErrNaoEFaxineiroDaSolicitacao
	}

	solicitacao, err := s.solicitacoes.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if err != nil {
		return nil, err
	}
	if solicitacao == nil {
		return nil, errosdominio.ErrSolicitacaoNaoEncontrada
	}

	// Verificar novamente se o horário ainda está disponível
	dataFim := solicitacao.DataFimEstimada(limpeza.DuracaoEstimada)
	if err := s.agenda.VerificarDisponibilidade(ctx, faxineiroID, solicitacao.DataAgendada, dataFim); err != nil {
		return nil, err
	}

	if err := solicitacao.Aceitar(); err != nil {
		return nil, err
	}

	// Bloquear o horário na agenda ANTES de salvar a solicitação como aceita.
	// Se o bloqueio falhar, a solicitação não é persistida como aceita.
	if _, err := s.agenda.CriarBloqueioServico(ctx, faxineiroID, solicitacao.ID, solicitacao.DataAgendada, dataFim); err != nil {
		return nil, err
	}

	if err := s.solicitacoes.Atualizar(ctx, solicitacao); err != nil {
		// Se falhar ao salvar a solicitação, liberar o bloqueio criado
		_ = s.agenda.LiberarBloqueioPorSolicitacao(ctx, solicitacao.ID)
		return nil, err
	}

	return solicitacao, nil
}

// RejeitarSolicitacao permite que o faxineiro rejeite a solicitação de um cliente.
func (s *ServicoSolicitacao) RejeitarSolicitacao(ctx context.Context, faxineiroID, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}
	if err := limpeza.VerificarPropriedade(faxineiroID); err != nil {
		return nil, errosdominio.ErrNaoEFaxineiroDaSolicitacao
	}

	solicitacao, err := s.solicitacoes.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if err != nil {
		return nil, err
	}
	if solicitacao == nil {
		return nil, errosdominio.ErrSolicitacaoNaoEncontrada
	}

	if err := solicitacao.Rejeitar(); err != nil {
		return nil, err
	}
	if err := s.solicitacoes.Atualizar(ctx, solicitacao); err != nil {
		return nil, err
	}

	return solicitacao, nil
}

// CancelarSolicitacao permite que o cliente cancele sua própria solicitação (pendente ou aceita).
// Se a solicitação estava aceita, libera o horário na agenda do faxineiro.
func (s *ServicoSolicitacao) CancelarSolicitacao(ctx context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	solicitacao, err := s.solicitacoes.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if err != nil {
		return nil, err
	}
	if solicitacao == nil {
		return nil, errosdominio.ErrSolicitacaoNaoEncontrada
	}
	if solicitacao.ClienteID != clienteID {
		return nil, errosdominio.ErrNaoEClienteSolicitante
	}

	estaAceita := solicitacao.Status == valueobject.StatusSolicitacaoAceita

	if err := solicitacao.Cancelar(time.Now()); err != nil {
		return nil, err
	}
	if err := s.solicitacoes.Atualizar(ctx, solicitacao); err != nil {
		return nil, err
	}

	// Se estava aceita, liberar o bloqueio na agenda
	if estaAceita {
		if err := s.agenda.LiberarBloqueioPorSolicitacao(ctx, solicitacao.ID); err != nil {
			return nil, err
		}
	}

	return solicitacao, nil
}

// ListarSolicitacoesPorLimpeza retorna todas as solicitações de um serviço. Apenas o faxineiro deve chamar.
func (s *ServicoSolicitacao) ListarSolicitacoesPorLimpeza(ctx context.Context, faxineiroID, limpezaID int) ([]*entity.Solicitacao, error) {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}
	if err := limpeza.VerificarPropriedade(faxineiroID); err != nil {
		return nil, errosdominio.ErrNaoEFaxineiroDaSolicitacao
	}
	return s.solicitacoes.ListarPorLimpeza(ctx, limpezaID)
}

// ListarSolicitacoesPorCliente retorna todas as solicitações feitas por um cliente.
func (s *ServicoSolicitacao) ListarSolicitacoesPorCliente(ctx context.Context, clienteID int) ([]*entity.Solicitacao, error) {
	return s.solicitacoes.ListarPorCliente(ctx, clienteID)
}
