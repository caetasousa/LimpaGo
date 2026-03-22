package service

import (
	"context"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/repository"
)

// ServicoAgenda gerencia a disponibilidade e os bloqueios de horário dos faxineiros.
type ServicoAgenda struct {
	agenda repository.RepositorioAgenda
}

func NovoServicoAgenda(agenda repository.RepositorioAgenda) *ServicoAgenda {
	return &ServicoAgenda{agenda: agenda}
}

// AdicionarDisponibilidade permite que o faxineiro defina um bloco de horário disponível.
func (s *ServicoAgenda) AdicionarDisponibilidade(ctx context.Context, faxineiroID int, diaSemana time.Weekday, horaInicio, horaFim int) (*entity.Disponibilidade, error) {
	d, err := entity.NovaDisponibilidade(faxineiroID, diaSemana, horaInicio, horaFim)
	if err != nil {
		return nil, err
	}

	if err := s.agenda.SalvarDisponibilidade(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// RemoverDisponibilidade remove um bloco de disponibilidade do faxineiro.
func (s *ServicoAgenda) RemoverDisponibilidade(ctx context.Context, id, faxineiroID int) error {
	return s.agenda.DeletarDisponibilidade(ctx, id, faxineiroID)
}

// ListarDisponibilidade retorna todos os blocos de disponibilidade de um faxineiro.
func (s *ServicoAgenda) ListarDisponibilidade(ctx context.Context, faxineiroID int) ([]*entity.Disponibilidade, error) {
	return s.agenda.ListarDisponibilidadePorFaxineiro(ctx, faxineiroID)
}

// VerificarDisponibilidade verifica se o faxineiro está disponível em um determinado período.
// Retorna erro se:
// 1. O horário não cai dentro de um bloco de disponibilidade semanal
// 2. Já existe um bloqueio (serviço agendado) que conflita com o período
func (s *ServicoAgenda) VerificarDisponibilidade(ctx context.Context, faxineiroID int, inicio, fim time.Time) error {
	// Verificar se o faxineiro tem disponibilidade nesse dia da semana e horário
	disponibilidades, err := s.agenda.ListarDisponibilidadePorDia(ctx, faxineiroID, inicio.Weekday())
	if err != nil {
		return err
	}

	horaInicio := inicio.Hour()
	horaFim := fim.Hour()
	if fim.Minute() > 0 {
		horaFim++
	}

	disponivel := false
	for _, d := range disponibilidades {
		if d.HoraInicio <= horaInicio && d.HoraFim >= horaFim {
			disponivel = true
			break
		}
	}
	if !disponivel {
		return errosdominio.ErrHorarioIndisponivel
	}

	// Verificar se não há conflito com bloqueios existentes
	bloqueios, err := s.agenda.ListarBloqueiosPorPeriodo(ctx, faxineiroID, inicio, fim)
	if err != nil {
		return err
	}
	if len(bloqueios) > 0 {
		return errosdominio.ErrConflitoAgenda
	}

	return nil
}

// CriarBloqueioServico reserva um horário na agenda do faxineiro para uma solicitação aceita.
func (s *ServicoAgenda) CriarBloqueioServico(ctx context.Context, faxineiroID, solicitacaoID int, inicio, fim time.Time) (*entity.Bloqueio, error) {
	bloqueio, err := entity.NovoBloqueioServico(faxineiroID, solicitacaoID, inicio, fim)
	if err != nil {
		return nil, err
	}

	if err := s.agenda.SalvarBloqueio(ctx, bloqueio); err != nil {
		return nil, err
	}
	return bloqueio, nil
}

// CriarBloqueioPessoal permite que o faxineiro bloqueie um horário pessoal (ex: consulta, folga).
func (s *ServicoAgenda) CriarBloqueioPessoal(ctx context.Context, faxineiroID int, inicio, fim time.Time) (*entity.Bloqueio, error) {
	bloqueio, err := entity.NovoBloqueiopessoal(faxineiroID, inicio, fim)
	if err != nil {
		return nil, err
	}

	if err := s.agenda.SalvarBloqueio(ctx, bloqueio); err != nil {
		return nil, err
	}
	return bloqueio, nil
}

// RemoverBloqueioPessoal remove um bloqueio pessoal do faxineiro.
// Apenas bloqueios pessoais podem ser removidos por esta função.
func (s *ServicoAgenda) RemoverBloqueioPessoal(ctx context.Context, id, faxineiroID int) error {
	bloqueio, err := s.agenda.BuscarBloqueioPorID(ctx, id)
	if err != nil {
		return err
	}
	if bloqueio == nil {
		return errosdominio.ErrBloqueioNaoEncontrado
	}
	if bloqueio.FaxineiroID != faxineiroID {
		return errosdominio.ErrNaoEFaxineiroDoBloqueio
	}
	if !bloqueio.EPessoal() {
		return errosdominio.ErrBloqueioPessoalApenas
	}
	return s.agenda.DeletarBloqueio(ctx, id)
}

// ListarBloqueios retorna todos os bloqueios (serviço e pessoal) de um faxineiro.
func (s *ServicoAgenda) ListarBloqueios(ctx context.Context, faxineiroID int) ([]*entity.Bloqueio, error) {
	return s.agenda.ListarBloqueiosPorFaxineiro(ctx, faxineiroID)
}

// LiberarBloqueioPorSolicitacao remove o bloqueio associado a uma solicitação (ex: quando cancelada).
func (s *ServicoAgenda) LiberarBloqueioPorSolicitacao(ctx context.Context, solicitacaoID int) error {
	bloqueio, err := s.agenda.BuscarBloqueioPorSolicitacao(ctx, solicitacaoID)
	if err != nil {
		return err
	}
	if bloqueio == nil {
		return nil // nenhum bloqueio a liberar
	}
	return s.agenda.DeletarBloqueio(ctx, bloqueio.ID)
}
