package testutil

import (
	"context"

	"phresh-go/domain/entity"
	"phresh-go/domain/valueobject"
)

type chaveClienteLimpeza struct {
	ClienteID int
	LimpezaID int
}

type RepositorioSolicitacaoMock struct {
	solicitacoes map[chaveClienteLimpeza]*entity.Solicitacao
	proximoID    int
}

func NovoRepositorioSolicitacaoMock() *RepositorioSolicitacaoMock {
	return &RepositorioSolicitacaoMock{
		solicitacoes: make(map[chaveClienteLimpeza]*entity.Solicitacao),
		proximoID:    1,
	}
}

func (r *RepositorioSolicitacaoMock) BuscarPorClienteELimpeza(_ context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	chave := chaveClienteLimpeza{ClienteID: clienteID, LimpezaID: limpezaID}
	return r.solicitacoes[chave], nil
}

func (r *RepositorioSolicitacaoMock) BuscarAtivaPorClienteELimpeza(_ context.Context, clienteID, limpezaID int) (*entity.Solicitacao, error) {
	chave := chaveClienteLimpeza{ClienteID: clienteID, LimpezaID: limpezaID}
	s := r.solicitacoes[chave]
	if s == nil {
		return nil, nil
	}
	if s.Status == valueobject.StatusSolicitacaoPendente || s.Status == valueobject.StatusSolicitacaoAceita {
		return s, nil
	}
	return nil, nil
}

func (r *RepositorioSolicitacaoMock) ListarPorLimpeza(_ context.Context, limpezaID int) ([]*entity.Solicitacao, error) {
	var resultado []*entity.Solicitacao
	for _, s := range r.solicitacoes {
		if s.LimpezaID == limpezaID {
			resultado = append(resultado, s)
		}
	}
	return resultado, nil
}

func (r *RepositorioSolicitacaoMock) ListarPorCliente(_ context.Context, clienteID int) ([]*entity.Solicitacao, error) {
	var resultado []*entity.Solicitacao
	for _, s := range r.solicitacoes {
		if s.ClienteID == clienteID {
			resultado = append(resultado, s)
		}
	}
	return resultado, nil
}

func (r *RepositorioSolicitacaoMock) Salvar(_ context.Context, solicitacao *entity.Solicitacao) error {
	solicitacao.ID = r.proximoID
	r.proximoID++
	chave := chaveClienteLimpeza{ClienteID: solicitacao.ClienteID, LimpezaID: solicitacao.LimpezaID}
	r.solicitacoes[chave] = solicitacao
	return nil
}

func (r *RepositorioSolicitacaoMock) Atualizar(_ context.Context, solicitacao *entity.Solicitacao) error {
	chave := chaveClienteLimpeza{ClienteID: solicitacao.ClienteID, LimpezaID: solicitacao.LimpezaID}
	r.solicitacoes[chave] = solicitacao
	return nil
}

func (r *RepositorioSolicitacaoMock) Deletar(_ context.Context, clienteID, limpezaID int) error {
	chave := chaveClienteLimpeza{ClienteID: clienteID, LimpezaID: limpezaID}
	delete(r.solicitacoes, chave)
	return nil
}
