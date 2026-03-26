package postgres

import (
	"context"
	"database/sql"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

// RepositorioSolicitacaoPG implementa repository.RepositorioSolicitacao com PostgreSQL.
type RepositorioSolicitacaoPG struct {
	db *sql.DB
}

// NovoRepositorioSolicitacaoPG cria um novo RepositorioSolicitacaoPG.
func NovoRepositorioSolicitacaoPG(db *sql.DB) *RepositorioSolicitacaoPG {
	return &RepositorioSolicitacaoPG{db: db}
}

const colunassolicitacao = `id, cliente_id, limpeza_id, status, data_agendada, preco_total, multa_cancelamento,
	endereco_rua, endereco_complemento, endereco_bairro, endereco_cidade, endereco_estado, endereco_cep,
	criado_em, atualizado_em`

func (r *RepositorioSolicitacaoPG) BuscarPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	q := `SELECT ` + colunassolicitacao + ` FROM solicitacoes WHERE cliente_id=$1 AND limpeza_id=$2 ORDER BY criado_em DESC LIMIT 1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, clienteID, limpezaID))
}

func (r *RepositorioSolicitacaoPG) BuscarAtivaPorClienteELimpeza(ctx context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	q := `SELECT ` + colunassolicitacao + ` FROM solicitacoes
	      WHERE cliente_id=$1 AND limpeza_id=$2 AND status IN ('pendente','aceita') LIMIT 1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, clienteID, limpezaID))
}

func (r *RepositorioSolicitacaoPG) ListarPorLimpeza(ctx context.Context, limpezaID int) ([]*entity.Solicitacao, error) {
	q := `SELECT ` + colunassolicitacao + ` FROM solicitacoes WHERE limpeza_id=$1 ORDER BY criado_em DESC`
	return r.escanearLista(obterExecutor(ctx, r.db).QueryContext(ctx, q, limpezaID))
}

func (r *RepositorioSolicitacaoPG) ListarPorCliente(ctx context.Context, clienteID int) ([]*entity.Solicitacao, error) {
	q := `SELECT ` + colunassolicitacao + ` FROM solicitacoes WHERE cliente_id=$1 ORDER BY criado_em DESC`
	return r.escanearLista(obterExecutor(ctx, r.db).QueryContext(ctx, q, clienteID))
}

func (r *RepositorioSolicitacaoPG) Salvar(ctx context.Context, s *entity.Solicitacao) error {
	q := `INSERT INTO solicitacoes
	      (cliente_id, limpeza_id, status, data_agendada, preco_total, multa_cancelamento,
	       endereco_rua, endereco_complemento, endereco_bairro, endereco_cidade, endereco_estado, endereco_cep)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	      RETURNING id, criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		s.ClienteID, s.LimpezaID, string(s.Status), s.DataAgendada, s.PrecoTotal, s.MultaCancelamento,
		s.Endereco.Rua, s.Endereco.Complemento, s.Endereco.Bairro,
		s.Endereco.Cidade, s.Endereco.Estado, s.Endereco.CEP).
		Scan(&s.ID, &s.CriadoEm, &s.AtualizadoEm)
}

func (r *RepositorioSolicitacaoPG) Atualizar(ctx context.Context, s *entity.Solicitacao) error {
	q := `UPDATE solicitacoes SET status=$1, multa_cancelamento=$2, atualizado_em=NOW()
	      WHERE id=$3 RETURNING atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q, string(s.Status), s.MultaCancelamento, s.ID).
		Scan(&s.AtualizadoEm)
}

func (r *RepositorioSolicitacaoPG) Deletar(ctx context.Context, clienteID, limpezaID int) error {
	_, err := obterExecutor(ctx, r.db).ExecContext(ctx,
		`DELETE FROM solicitacoes WHERE cliente_id=$1 AND limpeza_id=$2`, clienteID, limpezaID)
	return err
}

func (r *RepositorioSolicitacaoPG) escanearUm(row *sql.Row) (*entity.Solicitacao, error) {
	s := &entity.Solicitacao{}
	var status string
	err := row.Scan(
		&s.ID, &s.ClienteID, &s.LimpezaID, &status, &s.DataAgendada, &s.PrecoTotal, &s.MultaCancelamento,
		&s.Endereco.Rua, &s.Endereco.Complemento, &s.Endereco.Bairro,
		&s.Endereco.Cidade, &s.Endereco.Estado, &s.Endereco.CEP,
		&s.CriadoEm, &s.AtualizadoEm,
	)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrSolicitacaoNaoEncontrada)
	}
	s.Status = valueobject.StatusSolicitacao(status)
	return s, nil
}

func (r *RepositorioSolicitacaoPG) escanearLista(rows *sql.Rows, err error) ([]*entity.Solicitacao, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*entity.Solicitacao
	for rows.Next() {
		s := &entity.Solicitacao{}
		var status string
		if err := rows.Scan(
			&s.ID, &s.ClienteID, &s.LimpezaID, &status, &s.DataAgendada, &s.PrecoTotal, &s.MultaCancelamento,
			&s.Endereco.Rua, &s.Endereco.Complemento, &s.Endereco.Bairro,
			&s.Endereco.Cidade, &s.Endereco.Estado, &s.Endereco.CEP,
			&s.CriadoEm, &s.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		s.Status = valueobject.StatusSolicitacao(status)
		lista = append(lista, s)
	}
	return lista, rows.Err()
}
