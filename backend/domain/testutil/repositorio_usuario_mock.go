package testutil

import (
	"context"

	"limpaGo/domain/entity"
)

type RepositorioUsuarioMock struct {
	usuarios   map[int]*entity.Usuario
	porEmail   map[string]*entity.Usuario
	porNome    map[string]*entity.Usuario
	proximoID  int
}

func NovoRepositorioUsuarioMock() *RepositorioUsuarioMock {
	return &RepositorioUsuarioMock{
		usuarios:  make(map[int]*entity.Usuario),
		porEmail:  make(map[string]*entity.Usuario),
		porNome:   make(map[string]*entity.Usuario),
		proximoID: 1,
	}
}

func (r *RepositorioUsuarioMock) BuscarPorID(_ context.Context, id int) (*entity.Usuario, error) {
	return r.usuarios[id], nil
}

func (r *RepositorioUsuarioMock) BuscarPorEmail(_ context.Context, email string) (*entity.Usuario, error) {
	return r.porEmail[email], nil
}

func (r *RepositorioUsuarioMock) BuscarPorNomeUsuario(_ context.Context, nomeUsuario string) (*entity.Usuario, error) {
	return r.porNome[nomeUsuario], nil
}

func (r *RepositorioUsuarioMock) Salvar(_ context.Context, usuario *entity.Usuario) error {
	usuario.ID = r.proximoID
	r.proximoID++
	r.usuarios[usuario.ID] = usuario
	r.porEmail[usuario.Email] = usuario
	r.porNome[usuario.NomeUsuario] = usuario
	return nil
}
