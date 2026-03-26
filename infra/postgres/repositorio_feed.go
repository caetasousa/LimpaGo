package postgres

import (
	"context"
	"database/sql"

	"limpaGo/domain/entity"
	"limpaGo/domain/valueobject"
)

// RepositorioFeedPG implementa repository.RepositorioFeed com PostgreSQL.
type RepositorioFeedPG struct {
	db *sql.DB
}

// NovoRepositorioFeedPG cria um novo RepositorioFeedPG.
func NovoRepositorioFeedPG(db *sql.DB) *RepositorioFeedPG {
	return &RepositorioFeedPG{db: db}
}

func (r *RepositorioFeedPG) BuscarPaginaFeed(ctx context.Context, pagina, tamanhoPagina int) (*entity.PaginaFeed, error) {
	offset := (pagina - 1) * tamanhoPagina

	q := `SELECT
	          id, nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, profissional_id,
	          criado_em, atualizado_em,
	          CASE WHEN criado_em = atualizado_em THEN 'criacao' ELSE 'atualizacao' END AS tipo_evento,
	          ROW_NUMBER() OVER (ORDER BY atualizado_em DESC) AS numero_linha,
	          COUNT(*) OVER () AS total_itens
	      FROM limpezas
	      ORDER BY atualizado_em DESC
	      LIMIT $1 OFFSET $2`

	rows, err := obterExecutor(ctx, r.db).QueryContext(ctx, q, tamanhoPagina, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pagina_ := &entity.PaginaFeed{
		Pagina:        pagina,
		TamanhoPagina: tamanhoPagina,
	}

	for rows.Next() {
		l := &entity.Limpeza{}
		var tipoLimpeza, tipoEvento string
		var numeroLinha, totalItens int

		if err := rows.Scan(
			&l.ID, &l.Nome, &l.Descricao, &l.ValorHora, &l.DuracaoEstimada, &tipoLimpeza, &l.ProfissionalID,
			&l.CriadoEm, &l.AtualizadoEm,
			&tipoEvento, &numeroLinha, &totalItens,
		); err != nil {
			return nil, err
		}
		l.TipoLimpeza = valueobject.TipoLimpeza(tipoLimpeza)

		pagina_.TotalItens = totalItens
		pagina_.Itens = append(pagina_.Itens, &entity.ItemFeed{
			Limpeza:     l,
			TipoEvento:  valueobject.TipoEventoFeed(tipoEvento),
			DataEvento:  l.AtualizadoEm,
			NumeroLinha: numeroLinha,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if pagina_.Itens == nil {
		pagina_.Itens = []*entity.ItemFeed{}
	}

	return pagina_, nil
}
