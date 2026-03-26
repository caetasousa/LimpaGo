package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	errosdominio "limpaGo/domain/errors"
)

const codigoViolacaoUnica = "23505"

// mapearErroPG converte erros do PostgreSQL em erros sentinela do domínio.
// erroNaoEncontrado é o sentinel a retornar quando sql.ErrNoRows for encontrado.
func mapearErroPG(err error, erroNaoEncontrado error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return erroNaoEncontrado
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == codigoViolacaoUnica {
		switch pgErr.ConstraintName {
		case "usuarios_email_key":
			return errosdominio.ErrEmailJaUtilizado
		case "usuarios_nome_usuario_key":
			return errosdominio.ErrNomeUsuarioJaUtilizado
		case "avaliacoes_cliente_id_limpeza_id_key":
			return errosdominio.ErrAvaliacaoDuplicada
		}
	}

	return err
}
