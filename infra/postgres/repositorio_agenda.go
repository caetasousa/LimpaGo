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

func (r *RepositorioAgendaPG) ListarDisponibilidadePorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Disponibilidade, error) {
	q := `SELECT id, faxineiro_id, dia_semana, hora_inicio, hora_fim, criado_em, atualizado_em
	      FROM disponibilidades WHERE faxineiro_id=$1 ORDER BY dia_semana, hora_inicio`
	return r.escanearDisponibilidades(obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID))
}

func (r *RepositorioAgendaPG) ListarDisponibilidadePorDia(ctx context.Context, faxineiroID int, diaSemana time.Weekday) ([]*entity.Disponibilidade, error) {
	q := `SELECT id, faxineiro_id, dia_semana, hora_inicio, hora_fim, criado_em, atualizado_em
	      FROM disponibilidades WHERE faxineiro_id=$1 AND dia_semana=$2 ORDER BY hora_inicio`
	return r.escanearDisponibilidades(obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID, int(diaSemana)))
}

func (r *RepositorioAgendaPG) SalvarDisponibilidade(ctx context.Context, d *entity.Disponibilidade) error {
	q := `INSERT INTO disponibilidades (faxineiro_id, dia_semana, hora_inicio, hora_fim)
	      VALUES ($1,$2,$3,$4)
	      RETURNING id, criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		d.FaxineiroID, int(d.DiaSemana), d.HoraInicio, d.HoraFim).
		Scan(&d.ID, &d.CriadoEm, &d.AtualizadoEm)
}

func (r *RepositorioAgendaPG) DeletarDisponibilidade(ctx context.Context, id, faxineiroID int) error {
	_, err := obterExecutor(ctx, r.db).ExecContext(ctx,
		`DELETE FROM disponibilidades WHERE id=$1 AND faxineiro_id=$2`, id, faxineiroID)
	return err
}

// --- Bloqueios ---

func (r *RepositorioAgendaPG) ListarBloqueiosPorPeriodo(ctx context.Context, faxineiroID int, inicio, fim time.Time) ([]*entity.Bloqueio, error) {
	q := `SELECT id, faxineiro_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE faxineiro_id=$1 AND data_inicio < $3 AND data_fim > $2`
	return r.escanearBloqueios(obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID, inicio, fim))
}

func (r *RepositorioAgendaPG) ListarBloqueiosPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Bloqueio, error) {
	q := `SELECT id, faxineiro_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE faxineiro_id=$1 ORDER BY data_inicio`
	return r.escanearBloqueios(obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID))
}

func (r *RepositorioAgendaPG) BuscarBloqueioPorSolicitacao(ctx context.Context, solicitacaoID int) (*entity.Bloqueio, error) {
	q := `SELECT id, faxineiro_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE solicitacao_id=$1`
	return r.escanearBloqueioUnico(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, solicitacaoID))
}

func (r *RepositorioAgendaPG) BuscarBloqueioPorID(ctx context.Context, id int) (*entity.Bloqueio, error) {
	q := `SELECT id, faxineiro_id, solicitacao_id, data_inicio, data_fim, criado_em
	      FROM bloqueios WHERE id=$1`
	return r.escanearBloqueioUnico(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, id))
}

func (r *RepositorioAgendaPG) SalvarBloqueio(ctx context.Context, b *entity.Bloqueio) error {
	q := `INSERT INTO bloqueios (faxineiro_id, solicitacao_id, data_inicio, data_fim)
	      VALUES ($1,$2,$3,$4) RETURNING id, criado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		b.FaxineiroID, b.SolicitacaoID, b.DataInicio, b.DataFim).
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
		if err := rows.Scan(&d.ID, &d.FaxineiroID, &diaSemana, &d.HoraInicio, &d.HoraFim, &d.CriadoEm, &d.AtualizadoEm); err != nil {
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
		if err := rows.Scan(&b.ID, &b.FaxineiroID, &b.SolicitacaoID, &b.DataInicio, &b.DataFim, &b.CriadoEm); err != nil {
			return nil, err
		}
		lista = append(lista, b)
	}
	return lista, rows.Err()
}

func (r *RepositorioAgendaPG) escanearBloqueioUnico(row *sql.Row) (*entity.Bloqueio, error) {
	b := &entity.Bloqueio{}
	err := row.Scan(&b.ID, &b.FaxineiroID, &b.SolicitacaoID, &b.DataInicio, &b.DataFim, &b.CriadoEm)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrBloqueioNaoEncontrado)
	}
	return b, nil
}
