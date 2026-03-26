package repository

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioPerfil interface {
	// Perfil base
	BuscarPorUsuarioID(ctx context.Context, usuarioID int) (*entity.Perfil, error)
	Salvar(ctx context.Context, perfil *entity.Perfil) error
	Atualizar(ctx context.Context, perfil *entity.Perfil) error

	// Perfil Profissional
	BuscarPerfilProfissional(ctx context.Context, usuarioID int) (*entity.PerfilProfissional, error)
	SalvarPerfilProfissional(ctx context.Context, perfil *entity.PerfilProfissional) error
	AtualizarPerfilProfissional(ctx context.Context, perfil *entity.PerfilProfissional) error

	// Perfil Cliente
	BuscarPerfilCliente(ctx context.Context, usuarioID int) (*entity.PerfilCliente, error)
	SalvarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error
	AtualizarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error
}
