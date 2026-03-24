package postgres

import (
	"context"
	"database/sql"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

// RepositorioLimpezaPG implementa repository.RepositorioLimpeza com PostgreSQL.
type RepositorioLimpezaPG struct {
	db *sql.DB
}

// NovoRepositorioLimpezaPG cria um novo RepositorioLimpezaPG.
func NovoRepositorioLimpezaPG(db *sql.DB) *RepositorioLimpezaPG {
	return &RepositorioLimpezaPG{db: db}
}

func (r *RepositorioLimpezaPG) BuscarPorID(ctx context.Context, id int) (*entity.Limpeza, error) {
	q := `SELECT id, nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, faxineiro_id, criado_em, atualizado_em
	      FROM limpezas WHERE id = $1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, id))
}

func (r *RepositorioLimpezaPG) ListarPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Limpeza, error) {
	q := `SELECT id, nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, faxineiro_id, criado_em, atualizado_em
	      FROM limpezas WHERE faxineiro_id = $1 ORDER BY criado_em DESC`
	return r.escanearLista(obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID))
}

func (r *RepositorioLimpezaPG) ListarTodas(ctx context.Context, pagina, tamanhoPagina int) ([]*entity.Limpeza, error) {
	offset := (pagina - 1) * tamanhoPagina
	q := `SELECT id, nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, faxineiro_id, criado_em, atualizado_em
	      FROM limpezas ORDER BY criado_em DESC LIMIT $1 OFFSET $2`
	return r.escanearLista(obterExecutor(ctx, r.db).QueryContext(ctx, q, tamanhoPagina, offset))
}

func (r *RepositorioLimpezaPG) Salvar(ctx context.Context, limpeza *entity.Limpeza) error {
	q := `INSERT INTO limpezas (nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, faxineiro_id)
	      VALUES ($1,$2,$3,$4,$5,$6)
	      RETURNING id, criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		limpeza.Nome, limpeza.Descricao, limpeza.ValorHora, limpeza.DuracaoEstimada,
		string(limpeza.TipoLimpeza), limpeza.FaxineiroID).
		Scan(&limpeza.ID, &limpeza.CriadoEm, &limpeza.AtualizadoEm)
}

func (r *RepositorioLimpezaPG) Atualizar(ctx context.Context, limpeza *entity.Limpeza) error {
	q := `UPDATE limpezas SET nome=$1, descricao=$2, valor_hora=$3, duracao_estimada=$4, tipo_limpeza=$5, atualizado_em=NOW()
	      WHERE id=$6 RETURNING atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		limpeza.Nome, limpeza.Descricao, limpeza.ValorHora, limpeza.DuracaoEstimada,
		string(limpeza.TipoLimpeza), limpeza.ID).
		Scan(&limpeza.AtualizadoEm)
}

func (r *RepositorioLimpezaPG) Deletar(ctx context.Context, id int) error {
	_, err := obterExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM limpezas WHERE id = $1`, id)
	return err
}

func (r *RepositorioLimpezaPG) escanearUm(row *sql.Row) (*entity.Limpeza, error) {
	l := &entity.Limpeza{}
	var tipoLimpeza string
	err := row.Scan(&l.ID, &l.Nome, &l.Descricao, &l.ValorHora, &l.DuracaoEstimada, &tipoLimpeza, &l.FaxineiroID, &l.CriadoEm, &l.AtualizadoEm)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrLimpezaNaoEncontrada)
	}
	l.TipoLimpeza = valueobject.TipoLimpeza(tipoLimpeza)
	return l, nil
}

func (r *RepositorioLimpezaPG) escanearLista(rows *sql.Rows, err error) ([]*entity.Limpeza, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*entity.Limpeza
	for rows.Next() {
		l := &entity.Limpeza{}
		var tipoLimpeza string
		if err := rows.Scan(&l.ID, &l.Nome, &l.Descricao, &l.ValorHora, &l.DuracaoEstimada, &tipoLimpeza, &l.FaxineiroID, &l.CriadoEm, &l.AtualizadoEm); err != nil {
			return nil, err
		}
		l.TipoLimpeza = valueobject.TipoLimpeza(tipoLimpeza)
		lista = append(lista, l)
	}
	return lista, rows.Err()
}
