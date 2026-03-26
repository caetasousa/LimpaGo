package postgres

import (
	"context"
	"database/sql"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
)

// RepositorioUsuarioPG implementa repository.RepositorioUsuario com PostgreSQL.
type RepositorioUsuarioPG struct {
	db *sql.DB
}

// NovoRepositorioUsuarioPG cria um novo RepositorioUsuarioPG.
func NovoRepositorioUsuarioPG(db *sql.DB) *RepositorioUsuarioPG {
	return &RepositorioUsuarioPG{db: db}
}

func (r *RepositorioUsuarioPG) BuscarPorID(ctx context.Context, id int) (*entity.Usuario, error) {
	q := `SELECT id, email, nome_usuario, email_verificado, ativo, super_usuario, criado_em, atualizado_em
	      FROM usuarios WHERE id = $1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, id))
}

func (r *RepositorioUsuarioPG) BuscarPorEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	q := `SELECT id, email, nome_usuario, email_verificado, ativo, super_usuario, criado_em, atualizado_em
	      FROM usuarios WHERE email = $1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, email))
}

func (r *RepositorioUsuarioPG) BuscarPorNomeUsuario(ctx context.Context, nomeUsuario string) (*entity.Usuario, error) {
	q := `SELECT id, email, nome_usuario, email_verificado, ativo, super_usuario, criado_em, atualizado_em
	      FROM usuarios WHERE nome_usuario = $1`
	return r.escanearUm(obterExecutor(ctx, r.db).QueryRowContext(ctx, q, nomeUsuario))
}

func (r *RepositorioUsuarioPG) Salvar(ctx context.Context, usuario *entity.Usuario) error {
	q := `INSERT INTO usuarios (email, nome_usuario, email_verificado, ativo, super_usuario)
	      VALUES ($1, $2, $3, $4, $5)
	      RETURNING id, criado_em, atualizado_em`

	row := obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		usuario.Email, usuario.NomeUsuario, usuario.EmailVerificado, usuario.Ativo, usuario.SuperUsuario)

	err := row.Scan(&usuario.ID, &usuario.CriadoEm, &usuario.AtualizadoEm)
	return mapearErroPG(err, errosdominio.ErrUsuarioNaoEncontrado)
}

func (r *RepositorioUsuarioPG) escanearUm(row *sql.Row) (*entity.Usuario, error) {
	u := &entity.Usuario{}
	err := row.Scan(&u.ID, &u.Email, &u.NomeUsuario, &u.EmailVerificado, &u.Ativo, &u.SuperUsuario, &u.CriadoEm, &u.AtualizadoEm)
	if err != nil {
		mapped := mapearErroPG(err, nil)
		return nil, mapped
	}
	return u, nil
}
