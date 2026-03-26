package postgres

import (
	"context"
	"database/sql"
	"time"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
)

// RepositorioAgendaPG implementa repository.RepositorioAgenda com PostgreSQL.
type RepositorioAgendaPG struct {
	db *sql.DB
}

// NovoRepositorioAgendaPG cria um novo RepositorioAgendaPG.
func NovoRepositorioAgendaPG(db *sql.DB) *RepositorioAgendaPG {
	return &RepositorioAgendaPG{db: db}
}

// --- Disponibilidades ---

func (r *RepositorioAgendaPG) ListarDisponibilidadePorProfissional(ctx context.Context, profissionalID int) ([]*entity.Disponibilidade, error) {
	q := `SELECT id, profissional_id, dia_semana, hora_inicio, hora_fim, criado_em, atualizado_em
	      FROM disponibilidades WHERE profissional_id=$1 ORDER BY dia_semana, hora_inicio`
	return r.escanearDisponibilidades(obterExecutor(ctx, r.db).QueryContext(ctx, q, profissionalID))
}

func (r *RepositorioAgendaPG) ListarDisponibilidadePorDia(ctx context.Context, profissionalID int, diaSemana time.Weekday) ([]*entity.Disponibilidade, error) {
	q := `SELECT id, profissional_id, dia_semana, hora_inicio, hora_fim, criado_em, atualizado_em
	      FROM disponibilidades WHERE profissional_id=$1 AND dia_semana=$2 ORDER BY hora_inicio`
	return r.escanearDisponibilidades(obterExecutor(ctx, r.db).QueryContext(ctx, q, profissionalID, int(diaSemana)))
}

func (r *RepositorioAgendaPG) SalvarDisponibilidade(ctx context.Context, d *entity.Disponibilidade) error {
	q := `INSERT INTO disponibilidades (profissional_id, dia_semana, hora_inicio, hora_fim)
	      VALUES ($1,$2,$3,$4)
	      RETURNING id, criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		d.ProfissionalID, int(d.DiaSemana), d.HoraInicio, d.HoraFim).
		Scan(&d.ID, &d.CriadoEm, &d.AtualizadoEm)
}

func (r *RepositorioAgendaPG) DeletarDisponibilidade(ctx context.Context, id, profissionalID int) error {
	_, err := obterExecutor(ctx, r.db).ExecContext(ctx,
		`DELETE FROM disponibilidades WHERE id=$1 AND profissional_id=$2`, id, profissionalID)
	return err
}

// --- Bloqueios ---

func (r *RepositorioAgendaPG) ListarBloqueiosPorPeriodo(ctx context.Context, profissionalID int, inicio, fim time.Time) ([]*entity.Bloqueio, error) {
	q := `SELECT id, profissional_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE profissional_id=$1 AND data_inicio < $3 AND data_fim > $2`
	return r.escanearBloqueios(obterExecutor(ctx, r.db).QueryContext(ctx, q, profissionalID, inicio, fim))
}

func (r *RepositorioAgendaPG) ListarBloqueiosPorProfissional(ctx context.Context, profissionalID int) ([]*entity.Bloqueio, error) {
	q := `SELECT id, profissional_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE profissional_id=$1 ORDER BY data_inicio`
	return r.escanearBloqueios(obterExecutor(ctx, r.db).QueryContext(ctx, q, profissionalID))
}

func (r *RepositorioAgendaPG) BuscarBloqueioPorSolicitacao(ctx context.Context, solicitacaoID int) (*entity.Bloqueio, error) {
	q := `SELECT id, profissional_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE solicitacao_id=$1`
	return r.escanearBloqueioUnico(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, solicitacaoID))
}

func (r *RepositorioAgendaPG) BuscarBloqueioPorID(ctx context.Context, id int) (*entity.Bloqueio, error) {
	q := `SELECT id, profissional_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE id=$1`
	return r.escanearBloqueioUnico(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, id))
}

func (r *RepositorioAgendaPG) SalvarBloqueio(ctx context.Context, b *entity.Bloqueio) error {
	q := `INSERT INTO bloqueios (profissional_id, solicitacao_id, data_inicio, data_fim)
	      VALUES ($1,$2,$3,$4) RETURNING id, criado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		b.ProfissionalID, b.SolicitacaoID, b.DataInicio, b.DataFim).
		Scan(&b.ID, &b.CriadoEm)
}

func (r *RepositorioAgendaPG) DeletarBloqueio(ctx context.Context, id int) error {
	_, err := obterExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM bloqueios WHERE id=$1`, id)
	return err
}

func (r *RepositorioAgendaPG) escanearDisponibilidades(rows *sql.Rows, err error) ([]*entity.Disponibilidade, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*entity.Disponibilidade
	for rows.Next() {
		d := &entity.Disponibilidade{}
		var diaSemana int
		if err := rows.Scan(&d.ID, &d.ProfissionalID, &diaSemana, &d.HoraInicio, &d.HoraFim, &d.CriadoEm, &d.AtualizadoEm); err != nil {
			return nil, err
		}
		d.DiaSemana = time.Weekday(diaSemana)
		lista = append(lista, d)
	}
	return lista, rows.Err()
}

func (r *RepositorioAgendaPG) escanearBloqueios(rows *sql.Rows, err error) ([]*entity.Bloqueio, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*entity.Bloqueio
	for rows.Next() {
		b := &entity.Bloqueio{}
		if err := rows.Scan(&b.ID, &b.ProfissionalID, &b.SolicitacaoID, &b.DataInicio, &b.DataFim, &b.CriadoEm); err != nil {
			return nil, err
		}
		lista = append(lista, b)
	}
	return lista, rows.Err()
}

func (r *RepositorioAgendaPG) escanearBloqueioUnico(row *sql.Row) (*entity.Bloqueio, error) {
	b := &entity.Bloqueio{}
	err := row.Scan(&b.ID, &b.ProfissionalID, &b.SolicitacaoID, &b.DataInicio, &b.DataFim, &b.CriadoEm)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrBloqueioNaoEncontrado)
	}
	return b, nil
}
