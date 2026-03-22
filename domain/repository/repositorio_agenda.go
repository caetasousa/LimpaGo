package repository

import (
	"context"
	"time"

	"phresh-go/domain/entity"
)

type RepositorioAgenda interface {
	// Disponibilidade
	ListarDisponibilidadePorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Disponibilidade, error)
	ListarDisponibilidadePorDia(ctx context.Context, faxineiroID int, diaSemana time.Weekday) ([]*entity.Disponibilidade, error)
	SalvarDisponibilidade(ctx context.Context, d *entity.Disponibilidade) error
	DeletarDisponibilidade(ctx context.Context, id, faxineiroID int) error

	// Bloqueios (horários reservados — serviço ou pessoal)
	ListarBloqueiosPorPeriodo(ctx context.Context, faxineiroID int, inicio, fim time.Time) ([]*entity.Bloqueio, error)
	ListarBloqueiosPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Bloqueio, error)
	BuscarBloqueioPorSolicitacao(ctx context.Context, solicitacaoID int) (*entity.Bloqueio, error)
	BuscarBloqueioPorID(ctx context.Context, id int) (*entity.Bloqueio, error)
	SalvarBloqueio(ctx context.Context, b *entity.Bloqueio) error
	DeletarBloqueio(ctx context.Context, id int) error
}
