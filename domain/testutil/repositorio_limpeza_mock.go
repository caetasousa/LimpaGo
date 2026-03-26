package testutil

import (
	"context"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
)

type RepositorioLimpezaMock struct {
	limpezas  map[int]*entity.Limpeza
	proximoID int
}

func NovoRepositorioLimpezaMock() *RepositorioLimpezaMock {
	return &RepositorioLimpezaMock{
		limpezas:  make(map[int]*entity.Limpeza),
		proximoID: 1,
	}
}

func (r *RepositorioLimpezaMock) BuscarPorID(_ context.Context, id int) (*entity.Limpeza, error) {
	l, ok := r.limpezas[id]
	if !ok {
		return nil, errosdominio.ErrLimpezaNaoEncontrada
	}
	return l, nil
}

func (r *RepositorioLimpezaMock) ListarPorProfissional(_ context.Context, profissionalID int) ([]*entity.Limpeza, error) {
	var resultado []*entity.Limpeza
	for _, l := range r.limpezas {
		if l.ProfissionalID == profissionalID {
			resultado = append(resultado, l)
		}
	}
	return resultado, nil
}

func (r *RepositorioLimpezaMock) ListarTodas(_ context.Context, pagina, tamanhoPagina int) ([]*entity.Limpeza, error) {
	todas := make([]*entity.Limpeza, 0, len(r.limpezas))
	for _, l := range r.limpezas {
		todas = append(todas, l)
	}

	inicio := (pagina - 1) * tamanhoPagina
	if inicio >= len(todas) {
		return nil, nil
	}
	fim := inicio + tamanhoPagina
	if fim > len(todas) {
		fim = len(todas)
	}
	return todas[inicio:fim], nil
}

func (r *RepositorioLimpezaMock) Salvar(_ context.Context, limpeza *entity.Limpeza) error {
	limpeza.ID = r.proximoID
	r.proximoID++
	r.limpezas[limpeza.ID] = limpeza
	return nil
}

func (r *RepositorioLimpezaMock) Atualizar(_ context.Context, limpeza *entity.Limpeza) error {
	r.limpezas[limpeza.ID] = limpeza
	return nil
}

func (r *RepositorioLimpezaMock) Deletar(_ context.Context, id int) error {
	delete(r.limpezas, id)
	return nil
}
