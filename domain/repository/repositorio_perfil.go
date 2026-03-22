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

	// Perfil Faxineiro
	BuscarPerfilFaxineiro(ctx context.Context, usuarioID int) (*entity.PerfilFaxineiro, error)
	SalvarPerfilFaxineiro(ctx context.Context, perfil *entity.PerfilFaxineiro) error
	AtualizarPerfilFaxineiro(ctx context.Context, perfil *entity.PerfilFaxineiro) error

	// Perfil Cliente
	BuscarPerfilCliente(ctx context.Context, usuarioID int) (*entity.PerfilCliente, error)
	SalvarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error
	AtualizarPerfilCliente(ctx context.Context, perfil *entity.PerfilCliente) error
}
