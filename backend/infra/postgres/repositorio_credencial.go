package postgres

import (
	"context"
	"database/sql"

	"limpaGo/api/auth"
)

// RepositorioCredencialPG implementa auth.RepositorioCredencial com PostgreSQL.
type RepositorioCredencialPG struct {
	db *sql.DB
}

// NovoRepositorioCredencialPG cria um novo RepositorioCredencialPG.
func NovoRepositorioCredencialPG(db *sql.DB) *RepositorioCredencialPG {
	return &RepositorioCredencialPG{db: db}
}

func (r *RepositorioCredencialPG) BuscarPorUsuarioID(ctx context.Context, usuarioID int) (*auth.Credencial, error) {
	q := `SELECT usuario_id, senha_hash, criado_em, atualizado_em
	      FROM credenciais WHERE usuario_id = $1`

	c := &auth.Credencial{}
	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, usuarioID).
		Scan(&c.UsuarioID, &c.SenhaHash, &c.CriadoEm, &c.AtualizadoEm)
	if err != nil {
		return nil, mapearErroPG(err, nil)
	}
	return c, nil
}

func (r *RepositorioCredencialPG) Salvar(ctx context.Context, cred *auth.Credencial) error {
	q := `INSERT INTO credenciais (usuario_id, senha_hash)
	      VALUES ($1, $2)
	      ON CONFLICT (usuario_id) DO UPDATE SET senha_hash = EXCLUDED.senha_hash, atualizado_em = NOW()
	      RETURNING criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q, cred.UsuarioID, cred.SenhaHash).
		Scan(&cred.CriadoEm, &cred.AtualizadoEm)
}
