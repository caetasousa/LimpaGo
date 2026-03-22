package testutil

import (
	"context"

	"phresh-go/domain/entity"
)

type chaveAvaliacao struct {
	ClienteID int
	LimpezaID int
}

type RepositorioAvaliacaoMock struct {
	avaliacoes map[chaveAvaliacao]*entity.Avaliacao
	proximoID  int
}

func NovoRepositorioAvaliacaoMock() *RepositorioAvaliacaoMock {
	return &RepositorioAvaliacaoMock{
		avaliacoes: make(map[chaveAvaliacao]*entity.Avaliacao),
		proximoID:  1,
	}
}

func (r *RepositorioAvaliacaoMock) BuscarPorClienteELimpeza(_ context.Context, clienteID, limpezaID int) (*entity.Avaliacao, error) {
	chave := chaveAvaliacao{ClienteID: clienteID, LimpezaID: limpezaID}
	return r.avaliacoes[chave], nil
}

func (r *RepositorioAvaliacaoMock) ListarPorFaxineiro(_ context.Context, faxineiroID int) ([]*entity.Avaliacao, error) {
	var resultado []*entity.Avaliacao
	for _, a := range r.avaliacoes {
		if a.FaxineiroID == faxineiroID {
			resultado = append(resultado, a)
		}
	}
	return resultado, nil
}

func (r *RepositorioAvaliacaoMock) BuscarAgregadoPorFaxineiro(_ context.Context, faxineiroID int) (*entity.AgregadoAvaliacao, error) {
	var total int
	var soma float64
	for _, a := range r.avaliacoes {
		if a.FaxineiroID == faxineiroID {
			total++
			soma += float64(a.Nota)
		}
	}

	media := 0.0
	if total > 0 {
		media = soma / float64(total)
	}

	return &entity.AgregadoAvaliacao{
		FaxineiroID:     faxineiroID,
		MediaNota:       media,
		TotalAvaliacoes: total,
	}, nil
}

func (r *RepositorioAvaliacaoMock) Salvar(_ context.Context, avaliacao *entity.Avaliacao) error {
	avaliacao.ID = r.proximoID
	r.proximoID++
	chave := chaveAvaliacao{ClienteID: avaliacao.ClienteID, LimpezaID: avaliacao.LimpezaID}
	r.avaliacoes[chave] = avaliacao
	return nil
}
