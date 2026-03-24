package postgres

import (
	"context"
	"database/sql"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

// RepositorioAvaliacaoPG implementa repository.RepositorioAvaliacao com PostgreSQL.
type RepositorioAvaliacaoPG struct {
	db *sql.DB
}

// NovoRepositorioAvaliacaoPG cria um novo RepositorioAvaliacaoPG.
func NovoRepositorioAvaliacaoPG(db *sql.DB) *RepositorioAvaliacaoPG {
	return &RepositorioAvaliacaoPG{db: db}
}

func (r *RepositorioAvaliacaoPG) BuscarPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Avaliacao, error) {
	q := `SELECT id, limpeza_id, faxineiro_id, cliente_id, nota, comentario, criado_em
	      FROM avaliacoes WHERE cliente_id=$1 AND limpeza_id=$2`

	a := &entity.Avaliacao{}
	var nota int
	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, clienteID, limpezaID).
		Scan(&a.ID, &a.LimpezaID, &a.FaxineiroID, &a.ClienteID, &nota, &a.Comentario, &a.CriadoEm)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrAvaliacaoNaoEncontrada)
	}
	a.Nota = valueobject.Nota(nota)
	return a, nil
}

func (r *RepositorioAvaliacaoPG) ListarPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Avaliacao, error) {
	q := `SELECT id, limpeza_id, faxineiro_id, cliente_id, nota, comentario, criado_em
	      FROM avaliacoes WHERE faxineiro_id=$1 ORDER BY criado_em DESC`

	rows, err := obterExecutor(ctx, r.db).QueryContext(ctx, q, faxineiroID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*entity.Avaliacao
	for rows.Next() {
		a := &entity.Avaliacao{}
		var nota int
		if err := rows.Scan(&a.ID, &a.LimpezaID, &a.FaxineiroID, &a.ClienteID, &nota, &a.Comentario, &a.CriadoEm); err != nil {
			return nil, err
		}
		a.Nota = valueobject.Nota(nota)
		lista = append(lista, a)
	}
	return lista, rows.Err()
}

func (r *RepositorioAvaliacaoPG) BuscarAgregadoPorFaxineiro(ctx context.Context, faxineiroID int) (*entity.AgregadoAvaliacao, error) {
	q := `SELECT faxineiro_id, COALESCE(AVG(nota), 0)::float, COUNT(*)
	      FROM avaliacoes WHERE faxineiro_id=$1 GROUP BY faxineiro_id`

	ag := &entity.AgregadoAvaliacao{FaxineiroID: faxineiroID}
	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, faxineiroID).
		Scan(&ag.FaxineiroID, &ag.MediaNota, &ag.TotalAvaliacoes)
	if err != nil {
		if mapped := mapearErroPG(err, nil); mapped == nil {
			// Faxineiro sem avaliações — retorna agregado zerado
			return ag, nil
		}
		return nil, err
	}
	return ag, nil
}

func (r *RepositorioAvaliacaoPG) Salvar(ctx context.Context, avaliacao *entity.Avaliacao) error {
	q := `INSERT INTO avaliacoes (limpeza_id, faxineiro_id, cliente_id, nota, comentario)
	      VALUES ($1,$2,$3,$4,$5)
	      RETURNING id, criado_em`

	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		avaliacao.LimpezaID, avaliacao.FaxineiroID, avaliacao.ClienteID,
		int(avaliacao.Nota), avaliacao.Comentario).
		Scan(&avaliacao.ID, &avaliacao.CriadoEm)
	return mapearErroPG(err, nil)
}
