package testutil

import (
	"context"
	"time"

	"limpaGo/domain/entity"
)

type RepositorioAgendaMock struct {
	disponibilidades map[int]*entity.Disponibilidade
	bloqueios        map[int]*entity.Bloqueio
	proximoDispID    int
	proximoBloqID    int
}

func NovoRepositorioAgendaMock() *RepositorioAgendaMock {
	return &RepositorioAgendaMock{
		disponibilidades: make(map[int]*entity.Disponibilidade),
		bloqueios:        make(map[int]*entity.Bloqueio),
		proximoDispID:    1,
		proximoBloqID:    1,
	}
}

func (r *RepositorioAgendaMock) ListarDisponibilidadePorFaxineiro(_ context.Context, faxineiroID int) ([]*entity.Disponibilidade, error) {
	var resultado []*entity.Disponibilidade
	for _, d := range r.disponibilidades {
		if d.FaxineiroID == faxineiroID {
			resultado = append(resultado, d)
		}
	}
	return resultado, nil
}

func (r *RepositorioAgendaMock) ListarDisponibilidadePorDia(_ context.Context, faxineiroID int, diaSemana time.Weekday) ([]*entity.Disponibilidade, error) {
	var resultado []*entity.Disponibilidade
	for _, d := range r.disponibilidades {
		if d.FaxineiroID == faxineiroID && d.DiaSemana == diaSemana {
			resultado = append(resultado, d)
		}
	}
	return resultado, nil
}

func (r *RepositorioAgendaMock) SalvarDisponibilidade(_ context.Context, d *entity.Disponibilidade) error {
	d.ID = r.proximoDispID
	r.proximoDispID++
	r.disponibilidades[d.ID] = d
	return nil
}

func (r *RepositorioAgendaMock) DeletarDisponibilidade(_ context.Context, id, faxineiroID int) error {
	delete(r.disponibilidades, id)
	return nil
}

func (r *RepositorioAgendaMock) ListarBloqueiosPorPeriodo(_ context.Context, faxineiroID int, inicio, fim time.Time) ([]*entity.Bloqueio, error) {
	var resultado []*entity.Bloqueio
	for _, b := range r.bloqueios {
		if b.FaxineiroID == faxineiroID && b.DataInicio.Before(fim) && b.DataFim.After(inicio) {
			resultado = append(resultado, b)
		}
	}
	return resultado, nil
}

func (r *RepositorioAgendaMock) ListarBloqueiosPorFaxineiro(_ context.Context, faxineiroID int) ([]*entity.Bloqueio, error) {
	var resultado []*entity.Bloqueio
	for _, b := range r.bloqueios {
		if b.FaxineiroID == faxineiroID {
			resultado = append(resultado, b)
		}
	}
	return resultado, nil
}

func (r *RepositorioAgendaMock) BuscarBloqueioPorSolicitacao(_ context.Context, solicitacaoID int) (*entity.Bloqueio, error) {
	for _, b := range r.bloqueios {
		if b.SolicitacaoID != nil && *b.SolicitacaoID == solicitacaoID {
			return b, nil
		}
	}
	return nil, nil
}

func (r *RepositorioAgendaMock) BuscarBloqueioPorID(_ context.Context, id int) (*entity.Bloqueio, error) {
	return r.bloqueios[id], nil
}

func (r *RepositorioAgendaMock) SalvarBloqueio(_ context.Context, b *entity.Bloqueio) error {
	b.ID = r.proximoBloqID
	r.proximoBloqID++
	r.bloqueios[b.ID] = b
	return nil
}

func (r *RepositorioAgendaMock) DeletarBloqueio(_ context.Context, id int) error {
	delete(r.bloqueios, id)
	return nil
}
