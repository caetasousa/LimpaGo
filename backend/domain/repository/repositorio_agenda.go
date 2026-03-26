package repository

import (
	"context"
	"time"

	"limpaGo/domain/entity"
)

type RepositorioAgenda interface {
	// Disponibilidade
	ListarDisponibilidadePorProfissional(ctx context.Context, profissionalID int) ([]*entity.Disponibilidade, error)
	ListarDisponibilidadePorDia(ctx context.Context, profissionalID int, diaSemana time.Weekday) ([]*entity.Disponibilidade, error)
	SalvarDisponibilidade(ctx context.Context, d *entity.Disponibilidade) error
	DeletarDisponibilidade(ctx context.Context, id, profissionalID int) error

	// Bloqueios (horários reservados — serviço ou pessoal)
	ListarBloqueiosPorPeriodo(ctx context.Context, profissionalID int, inicio, fim time.Time) ([]*entity.Bloqueio, error)
	ListarBloqueiosPorProfissional(ctx context.Context, profissionalID int) ([]*entity.Bloqueio, error)
	BuscarBloqueioPorSolicitacao(ctx context.Context, solicitacaoID int) (*entity.Bloqueio, error)
	BuscarBloqueioPorID(ctx context.Context, id int) (*entity.Bloqueio, error)
	SalvarBloqueio(ctx context.Context, b *entity.Bloqueio) error
	DeletarBloqueio(ctx context.Context, id int) error
}
