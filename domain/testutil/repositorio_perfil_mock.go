package testutil

import (
	"context"

	"phresh-go/domain/entity"
)

type RepositorioPerfilMock struct {
	perfis           map[int]*entity.Perfil
	perfisFaxineiro  map[int]*entity.PerfilFaxineiro
	perfisCliente    map[int]*entity.PerfilCliente
}

func NovoRepositorioPerfilMock() *RepositorioPerfilMock {
	return &RepositorioPerfilMock{
		perfis:          make(map[int]*entity.Perfil),
		perfisFaxineiro: make(map[int]*entity.PerfilFaxineiro),
		perfisCliente:   make(map[int]*entity.PerfilCliente),
	}
}

func (r *RepositorioPerfilMock) BuscarPorUsuarioID(_ context.Context, usuarioID int) (*entity.Perfil, error) {
	return r.perfis[usuarioID], nil
}

func (r *RepositorioPerfilMock) Salvar(_ context.Context, perfil *entity.Perfil) error {
	r.perfis[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) Atualizar(_ context.Context, perfil *entity.Perfil) error {
	r.perfis[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) BuscarPerfilFaxineiro(_ context.Context, usuarioID int) (*entity.PerfilFaxineiro, error) {
	return r.perfisFaxineiro[usuarioID], nil
}

func (r *RepositorioPerfilMock) SalvarPerfilFaxineiro(_ context.Context, perfil *entity.PerfilFaxineiro) error {
	r.perfisFaxineiro[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) AtualizarPerfilFaxineiro(_ context.Context, perfil *entity.PerfilFaxineiro) error {
	r.perfisFaxineiro[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) BuscarPerfilCliente(_ context.Context, usuarioID int) (*entity.PerfilCliente, error) {
	return r.perfisCliente[usuarioID], nil
}

func (r *RepositorioPerfilMock) SalvarPerfilCliente(_ context.Context, perfil *entity.PerfilCliente) error {
	r.perfisCliente[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) AtualizarPerfilCliente(_ context.Context, perfil *entity.PerfilCliente) error {
	r.perfisCliente[perfil.UsuarioID] = perfil
	return nil
}
