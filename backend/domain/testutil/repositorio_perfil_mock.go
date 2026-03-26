package testutil

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioPerfilMock struct {
	perfis           map[int]*entity.Perfil
	perfisProfissional  map[int]*entity.PerfilProfissional
	perfisCliente    map[int]*entity.PerfilCliente
}

func NovoRepositorioPerfilMock() *RepositorioPerfilMock {
	return &RepositorioPerfilMock{
		perfis:          make(map[int]*entity.Perfil),
		perfisProfissional: make(map[int]*entity.PerfilProfissional),
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

func (r *RepositorioPerfilMock) BuscarPerfilProfissional(_ context.Context, usuarioID int) (*entity.PerfilProfissional, error) {
	return r.perfisProfissional[usuarioID], nil
}

func (r *RepositorioPerfilMock) SalvarPerfilProfissional(_ context.Context, perfil *entity.PerfilProfissional) error {
	r.perfisProfissional[perfil.UsuarioID] = perfil
	return nil
}

func (r *RepositorioPerfilMock) AtualizarPerfilProfissional(_ context.Context, perfil *entity.PerfilProfissional) error {
	r.perfisProfissional[perfil.UsuarioID] = perfil
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
