package service

import (
	"context"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/repository"
	"limpaGo/domain/valueobject"
)

// ServicoAvaliacao lida com notas e reputação de profissionais.
//
// Regras:
//   - Apenas o cliente que solicitou o serviço pode criar uma avaliação.
//   - A solicitação deve estar no estado aceita.
//   - Apenas uma avaliação por solicitação.
//   - Criar uma avaliação marca a solicitação como concluída.
type ServicoAvaliacao struct {
	avaliacoes   repository.RepositorioAvaliacao
	solicitacoes repository.RepositorioSolicitacao
	limpezas     repository.RepositorioLimpeza
}

func NovoServicoAvaliacao(
	avaliacoes repository.RepositorioAvaliacao,
	solicitacoes repository.RepositorioSolicitacao,
	limpezas repository.RepositorioLimpeza,
) *ServicoAvaliacao {
	return &ServicoAvaliacao{avaliacoes: avaliacoes, solicitacoes: solicitacoes, limpezas: limpezas}
}

func (s *ServicoAvaliacao) CriarAvaliacao(ctx context.Context, clienteID, limpezaID int, notaValor int, comentario string) (*entity.Avaliacao, error) {
	solicitacao, err := s.solicitacoes.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if err != nil {
		return nil, err
	}
	if solicitacao == nil {
		return nil, errosdominio.ErrSolicitacaoNaoEncontrada
	}
	if solicitacao.Status != valueobject.StatusSolicitacaoAceita {
		return nil, errosdominio.ErrSolicitacaoNaoAceita
	}

	existente, err := s.avaliacoes.BuscarPorClienteELimpeza(ctx, clienteID, limpezaID)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, errosdominio.ErrAvaliacaoDuplicada
	}

	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}

	nota, err := valueobject.NovaNota(notaValor)
	if err != nil {
		return nil, &entity.ErroValidacao{Campo: "nota", Mensagem: err.Error()}
	}

	avaliacao := entity.NovaAvaliacao(limpezaID, limpeza.ProfissionalID, clienteID, nota, comentario)

	if err := s.avaliacoes.Salvar(ctx, avaliacao); err != nil {
		return nil, err
	}

	solicitacao.Concluir()
	if err := s.solicitacoes.Atualizar(ctx, solicitacao); err != nil {
		return nil, err
	}

	return avaliacao, nil
}

func (s *ServicoAvaliacao) BuscarEstatisticasProfissional(ctx context.Context, profissionalID int) (*entity.AgregadoAvaliacao, error) {
	return s.avaliacoes.BuscarAgregadoPorProfissional(ctx, profissionalID)
}

func (s *ServicoAvaliacao) ListarAvaliacoesPorProfissional(ctx context.Context, profissionalID int) ([]*entity.Avaliacao, error) {
	return s.avaliacoes.ListarPorProfissional(ctx, profissionalID)
}
