package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

// RepositorioPerfilPG implementa repository.RepositorioPerfil com PostgreSQL.
type RepositorioPerfilPG struct {
	db *sql.DB
}

// NovoRepositorioPerfilPG cria um novo RepositorioPerfilPG.
func NovoRepositorioPerfilPG(db *sql.DB) *RepositorioPerfilPG {
	return &RepositorioPerfilPG{db: db}
}

// --- Perfil base ---

func (r *RepositorioPerfilPG) BuscarPorUsuarioID(ctx context.Context, usuarioID int) (*entity.Perfil, error) {
	q := `SELECT usuario_id, nome_completo, telefone, imagem, email, nome_usuario, criado_em, atualizado_em
	      FROM perfis WHERE usuario_id = $1`

	p := &entity.Perfil{}
	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, usuarioID).
		Scan(&p.UsuarioID, &p.NomeCompleto, &p.Telefone, &p.Imagem, &p.Email, &p.NomeUsuario, &p.CriadoEm, &p.AtualizadoEm)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrPerfilNaoEncontrado)
	}
	return p, nil
}

func (r *RepositorioPerfilPG) Salvar(ctx context.Context, perfil *entity.Perfil) error {
	q := `INSERT INTO perfis (usuario_id, nome_completo, telefone, imagem, email, nome_usuario)
	      VALUES ($1, $2, $3, $4, $5, $6)
	      RETURNING criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.UsuarioID, perfil.NomeCompleto, perfil.Telefone, perfil.Imagem, perfil.Email, perfil.NomeUsuario).
		Scan(&perfil.CriadoEm, &perfil.AtualizadoEm)
}

func (r *RepositorioPerfilPG) Atualizar(ctx context.Context, perfil *entity.Perfil) error {
	q := `UPDATE perfis SET nome_completo=$1, telefone=$2, imagem=$3, atualizado_em=NOW()
	      WHERE usuario_id=$4
	      RETURNING atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.NomeCompleto, perfil.Telefone, perfil.Imagem, perfil.UsuarioID).
		Scan(&perfil.AtualizadoEm)
}

// --- Perfil Faxineiro ---

func (r *RepositorioPerfilPG) BuscarPerfilFaxineiro(ctx context.Context, usuarioID int) (*entity.PerfilFaxineiro, error) {
	q := `SELECT usuario_id, descricao, anos_experiencia, especialidades, cidades_atendidas,
	             documento_rg, documento_cpf, foto_documento, verificado, criado_em, atualizado_em
	      FROM perfis_faxineiro WHERE usuario_id = $1`

	p := &entity.PerfilFaxineiro{}
	var especialidades pgtype.Array[string]
	var cidades pgtype.Array[string]

	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, usuarioID).Scan(
		&p.UsuarioID, &p.Descricao, &p.AnosExperiencia, &especialidades, &cidades,
		&p.DocumentoRG, &p.DocumentoCPF, &p.FotoDocumento, &p.Verificado,
		&p.CriadoEm, &p.AtualizadoEm,
	)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrPerfilFaxineiroNaoEncontrado)
	}
	p.Especialidades = especialidades.Elements
	p.CidadesAtendidas = cidades.Elements
	return p, nil
}

func (r *RepositorioPerfilPG) SalvarPerfilFaxineiro(ctx context.Context, perfil *entity.PerfilFaxineiro) error {
	q := `INSERT INTO perfis_faxineiro
	      (usuario_id, descricao, anos_experiencia, especialidades, cidades_atendidas,
	       documento_rg, documento_cpf, foto_documento, verificado)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	      RETURNING criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.UsuarioID, perfil.Descricao, perfil.AnosExperiencia,
		perfil.Especialidades, perfil.CidadesAtendidas,
		perfil.DocumentoRG, perfil.DocumentoCPF, perfil.FotoDocumento, perfil.Verificado).
		Scan(&perfil.CriadoEm, &perfil.AtualizadoEm)
}

func (r *RepositorioPerfilPG) AtualizarPerfilFaxineiro(ctx context.Context, perfil *entity.PerfilFaxineiro) error {
	q := `UPDATE perfis_faxineiro
	      SET descricao=$1, anos_experiencia=$2, especialidades=$3, cidades_atendidas=$4,
	          documento_rg=$5, documento_cpf=$6, foto_documento=$7, verificado=$8, atualizado_em=NOW()
	      WHERE usuario_id=$9
	      RETURNING atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.Descricao, perfil.AnosExperiencia,
		perfil.Especialidades, perfil.CidadesAtendidas,
		perfil.DocumentoRG, perfil.DocumentoCPF, perfil.FotoDocumento, perfil.Verificado,
		perfil.UsuarioID).
		Scan(&perfil.AtualizadoEm)
}

// --- Perfil Cliente ---

func (r *RepositorioPerfilPG) BuscarPerfilCliente(ctx context.Context, usuarioID int) (*entity.PerfilCliente, error) {
	q := `SELECT usuario_id,
	             endereco_rua, endereco_complemento, endereco_bairro,
	             endereco_cidade, endereco_estado, endereco_cep,
	             tipo_imovel, quartos, banheiros, tamanho_imovel_m2,
	             observacoes, faxineiro_preferido_id, criado_em, atualizado_em
	      FROM perfis_cliente WHERE usuario_id = $1`

	p := &entity.PerfilCliente{}
	var tipoImovel string
	err := obterExecutor(ctx, r.db).QueryRowContext(ctx, q, usuarioID).Scan(
		&p.UsuarioID,
		&p.Endereco.Rua, &p.Endereco.Complemento, &p.Endereco.Bairro,
		&p.Endereco.Cidade, &p.Endereco.Estado, &p.Endereco.CEP,
		&tipoImovel, &p.Quartos, &p.Banheiros, &p.TamanhoImovelM2,
		&p.Observacoes, &p.FaxineiroPreferidoID, &p.CriadoEm, &p.AtualizadoEm,
	)
	if err != nil {
		return nil, mapearErroPG(err, errosdominio.ErrPerfilClienteNaoEncontrado)
	}
	p.TipoImovel = valueobject.TipoImovel(tipoImovel)
	return p, nil
}

func (r *RepositorioPerfilPG) SalvarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error {
	q := `INSERT INTO perfis_cliente
	      (usuario_id, endereco_rua, endereco_complemento, endereco_bairro,
	       endereco_cidade, endereco_estado, endereco_cep,
	       tipo_imovel, quartos, banheiros, tamanho_imovel_m2,
	       observacoes, faxineiro_preferido_id)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	      RETURNING criado_em, atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.UsuarioID,
		perfil.Endereco.Rua, perfil.Endereco.Complemento, perfil.Endereco.Bairro,
		perfil.Endereco.Cidade, perfil.Endereco.Estado, perfil.Endereco.CEP,
		string(perfil.TipoImovel), perfil.Quartos, perfil.Banheiros, perfil.TamanhoImovelM2,
		perfil.Observacoes, perfil.FaxineiroPreferidoID).
		Scan(&perfil.CriadoEm, &perfil.AtualizadoEm)
}

func (r *RepositorioPerfilPG) AtualizarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error {
	q := `UPDATE perfis_cliente
	      SET endereco_rua=$1, endereco_complemento=$2, endereco_bairro=$3,
	          endereco_cidade=$4, endereco_estado=$5, endereco_cep=$6,
	          tipo_imovel=$7, quartos=$8, banheiros=$9, tamanho_imovel_m2=$10,
	          observacoes=$11, faxineiro_preferido_id=$12, atualizado_em=NOW()
	      WHERE usuario_id=$13
	      RETURNING atualizado_em`

	return obterExecutor(ctx, r.db).QueryRowContext(ctx, q,
		perfil.Endereco.Rua, perfil.Endereco.Complemento, perfil.Endereco.Bairro,
		perfil.Endereco.Cidade, perfil.Endereco.Estado, perfil.Endereco.CEP,
		string(perfil.TipoImovel), perfil.Quartos, perfil.Banheiros, perfil.TamanhoImovelM2,
		perfil.Observacoes, perfil.FaxineiroPreferidoID,
		perfil.UsuarioID).
		Scan(&perfil.AtualizadoEm)
}
