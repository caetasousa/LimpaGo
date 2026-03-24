package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type chaveContextoTx struct{}

// Executor é implementado tanto por *sql.DB quanto por *sql.Tx.
type Executor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// ComTransacao executa fn dentro de uma transação. Faz commit se fn retornar nil,
// ou rollback em caso de erro. A transação é propagada via contexto.
func ComTransacao(ctx context.Context, db *sql.DB, fn func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transação: %w", err)
	}

	ctxComTx := context.WithValue(ctx, chaveContextoTx{}, tx)

	if err := fn(ctxComTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar transação: %w", err)
	}
	return nil
}

// obterExecutor retorna a transação do contexto (se houver) ou o pool de conexões.
func obterExecutor(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := ctx.Value(chaveContextoTx{}).(*sql.Tx); ok && tx != nil {
		return tx
	}
	return db
}
