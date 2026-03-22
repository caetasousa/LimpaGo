package dto

import (
	"errors"
	"net/http"

	"phresh-go/domain/entity"
	errosdominio "phresh-go/domain/errors"
)

// RespostaErro é a estrutura padrão de resposta de erro da API.
type RespostaErro struct {
	Codigo    int               `json:"codigo"`
	Mensagem  string            `json:"mensagem"`
	Detalhes  map[string]string `json:"detalhes,omitempty"`
}

// NovaRespostaErro cria uma resposta de erro simples.
func NovaRespostaErro(codigo int, mensagem string) RespostaErro {
	return RespostaErro{Codigo: codigo, Mensagem: mensagem}
}

// NovaRespostaErroValidacao cria uma resposta de erro com detalhes de campo.
func NovaRespostaErroValidacao(campo, mensagem string) RespostaErro {
	return RespostaErro{
		Codigo:   http.StatusUnprocessableEntity,
		Mensagem: "erro de validação",
		Detalhes: map[string]string{campo: mensagem},
	}
}

// MapearErroDominio converte um erro de domínio para o status HTTP e RespostaErro correspondentes.
func MapearErroDominio(err error) (int, RespostaErro) {
	// Erro de validação de campo
	var erroValidacao *entity.ErroValidacao
	if errors.As(err, &erroValidacao) {
		return http.StatusUnprocessableEntity, NovaRespostaErroValidacao(erroValidacao.Campo, erroValidacao.Mensagem)
	}

	// 404 — não encontrado
	switch {
	case errors.Is(err, errosdominio.ErrUsuarioNaoEncontrado),
		errors.Is(err, errosdominio.ErrPerfilNaoEncontrado),
		errors.Is(err, errosdominio.ErrPerfilFaxineiroNaoEncontrado),
		errors.Is(err, errosdominio.ErrPerfilClienteNaoEncontrado),
		errors.Is(err, errosdominio.ErrLimpezaNaoEncontrada),
		errors.Is(err, errosdominio.ErrSolicitacaoNaoEncontrada),
		errors.Is(err, errosdominio.ErrDisponibilidadeNaoEncontrada),
		errors.Is(err, errosdominio.ErrBloqueioNaoEncontrado),
		errors.Is(err, errosdominio.ErrAvaliacaoNaoEncontrada):
		return http.StatusNotFound, NovaRespostaErro(http.StatusNotFound, err.Error())

	// 409 — conflito / duplicado
	case errors.Is(err, errosdominio.ErrEmailJaUtilizado),
		errors.Is(err, errosdominio.ErrNomeUsuarioJaUtilizado),
		errors.Is(err, errosdominio.ErrPerfilFaxineiroJaExiste),
		errors.Is(err, errosdominio.ErrPerfilClienteJaExiste),
		errors.Is(err, errosdominio.ErrSolicitacaoDuplicada),
		errors.Is(err, errosdominio.ErrAvaliacaoDuplicada):
		return http.StatusConflict, NovaRespostaErro(http.StatusConflict, err.Error())

	// 403 — sem permissão
	case errors.Is(err, errosdominio.ErrNaoEFaxineiroDaLimpeza),
		errors.Is(err, errosdominio.ErrNaoEClienteSolicitante),
		errors.Is(err, errosdominio.ErrNaoEFaxineiroDaSolicitacao),
		errors.Is(err, errosdominio.ErrNaoEFaxineiroDoBloqueio):
		return http.StatusForbidden, NovaRespostaErro(http.StatusForbidden, err.Error())

	// 422 — regra de negócio violada
	case errors.Is(err, errosdominio.ErrFaxineiroNaoPodeSolicitarProprio),
		errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerAceita),
		errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerCancelada),
		errors.Is(err, errosdominio.ErrSolicitacaoNaoPodeSerRejeitada),
		errors.Is(err, errosdominio.ErrAgendamentoNoPassado),
		errors.Is(err, errosdominio.ErrHorarioIndisponivel),
		errors.Is(err, errosdominio.ErrConflitoAgenda),
		errors.Is(err, errosdominio.ErrBloqueioPessoalApenas),
		errors.Is(err, errosdominio.ErrSolicitacaoNaoAceita):
		return http.StatusUnprocessableEntity, NovaRespostaErro(http.StatusUnprocessableEntity, err.Error())
	}

	// 500 — erro interno
	return http.StatusInternalServerError, NovaRespostaErro(http.StatusInternalServerError, "erro interno do servidor")
}
