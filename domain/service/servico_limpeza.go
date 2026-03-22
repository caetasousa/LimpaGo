package service

import (
	"context"

	"phresh-go/domain/entity"
	"phresh-go/domain/repository"
	"phresh-go/domain/valueobject"
)

type ServicoLimpeza struct {
	limpezas repository.RepositorioLimpeza
}

func NovoServicoLimpeza(limpezas repository.RepositorioLimpeza) *ServicoLimpeza {
	return &ServicoLimpeza{limpezas: limpezas}
}

// Criar permite que um faxineiro publique um novo serviço de limpeza com valor por hora e duração estimada.
func (s *ServicoLimpeza) Criar(ctx context.Context, faxineiroID int, nome, descricao string, valorHora, duracaoEstimada float64, tipoLimpeza valueobject.TipoLimpeza) (*entity.Limpeza, error) {
	limpeza, err := entity.NovaLimpeza(faxineiroID, nome, valorHora, duracaoEstimada, tipoLimpeza)
	if err != nil {
		return nil, err
	}
	limpeza.Descricao = descricao

	if err := s.limpezas.Salvar(ctx, limpeza); err != nil {
		return nil, err
	}
	return limpeza, nil
}

// Atualizar permite que o faxineiro atualize seu serviço publicado.
func (s *ServicoLimpeza) Atualizar(ctx context.Context, limpezaID, faxineiroID int, nome, descricao string, valorHora, duracaoEstimada float64, tipoLimpeza valueobject.TipoLimpeza) (*entity.Limpeza, error) {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return nil, err
	}

	if err := limpeza.VerificarPropriedade(faxineiroID); err != nil {
		return nil, err
	}

	if nome != "" {
		limpeza.Nome = nome
	}
	if descricao != "" {
		limpeza.Descricao = descricao
	}
	if valorHora > 0 {
		limpeza.ValorHora = valorHora
	}
	if duracaoEstimada > 0 {
		limpeza.DuracaoEstimada = duracaoEstimada
	}
	if tipoLimpeza != "" {
		if err := tipoLimpeza.Validar(); err == nil {
			limpeza.TipoLimpeza = tipoLimpeza
		}
	}

	if err := s.limpezas.Atualizar(ctx, limpeza); err != nil {
		return nil, err
	}
	return limpeza, nil
}

// Deletar permite que o faxineiro remova seu serviço publicado.
func (s *ServicoLimpeza) Deletar(ctx context.Context, limpezaID, faxineiroID int) error {
	limpeza, err := s.limpezas.BuscarPorID(ctx, limpezaID)
	if err != nil {
		return err
	}
	if err := limpeza.VerificarPropriedade(faxineiroID); err != nil {
		return err
	}
	return s.limpezas.Deletar(ctx, limpezaID)
}

func (s *ServicoLimpeza) BuscarPorID(ctx context.Context, id int) (*entity.Limpeza, error) {
	return s.limpezas.BuscarPorID(ctx, id)
}

// ListarPorFaxineiro retorna todos os serviços publicados por um faxineiro.
func (s *ServicoLimpeza) ListarPorFaxineiro(ctx context.Context, faxineiroID int) ([]*entity.Limpeza, error) {
	return s.limpezas.ListarPorFaxineiro(ctx, faxineiroID)
}

// ListarCatalogo retorna todos os serviços disponíveis para os clientes navegarem.
func (s *ServicoLimpeza) ListarCatalogo(ctx context.Context, pagina, tamanhoPagina int) ([]*entity.Limpeza, error) {
	p := valueobject.NovaPaginacao(pagina, tamanhoPagina)
	return s.limpezas.ListarTodas(ctx, p.Pagina, p.TamanhoPagina)
}
