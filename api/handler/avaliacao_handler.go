package handler

import (
	"net/http"

	"limpaGo/api/dto"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
)

// HandlerAvaliacao gerencia os endpoints de avaliações de faxineiros.
type HandlerAvaliacao struct {
	servico *service.ServicoAvaliacao
}

// NovoHandlerAvaliacao cria um novo HandlerAvaliacao.
func NovoHandlerAvaliacao(servico *service.ServicoAvaliacao) *HandlerAvaliacao {
	return &HandlerAvaliacao{servico: servico}
}

// CriarAvaliacao godoc
// @Summary Criar avaliação de um serviço concluído
// @Tags avaliacoes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoCriarAvaliacao true "Dados da avaliação"
// @Success 201 {object} dto.RespostaAvaliacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 409 {object} dto.RespostaErro
// @Failure 422 {object} dto.RespostaErro
// @Router /avaliacoes [post]
func (h *HandlerAvaliacao) CriarAvaliacao(w http.ResponseWriter, r *http.Request) {
	clienteID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoCriarAvaliacao
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	avaliacao, err := h.servico.CriarAvaliacao(r.Context(), clienteID, req.LimpezaID, req.Nota, req.Comentario)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DeAvaliacao(avaliacao))
}

// ListarAvaliacoes godoc
// @Summary Listar avaliações de um faxineiro
// @Tags avaliacoes
// @Produce json
// @Param faxineiro_id path int true "ID do faxineiro"
// @Success 200 {array} dto.RespostaAvaliacao
// @Failure 404 {object} dto.RespostaErro
// @Router /faxineiros/{faxineiro_id}/avaliacoes [get]
func (h *HandlerAvaliacao) ListarAvaliacoes(w http.ResponseWriter, r *http.Request) {
	faxineiroID, err := lerParamInteiro(r, "faxineiro_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "faxineiro_id inválido"))
		return
	}
	lista, err := h.servico.ListarAvaliacoesPorFaxineiro(r.Context(), faxineiroID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeAvaliacaoLista(lista))
}

// BuscarEstatisticas godoc
// @Summary Buscar estatísticas de avaliações de um faxineiro
// @Tags avaliacoes
// @Produce json
// @Param faxineiro_id path int true "ID do faxineiro"
// @Success 200 {object} dto.RespostaEstatisticasFaxineiro
// @Failure 404 {object} dto.RespostaErro
// @Router /faxineiros/{faxineiro_id}/estatisticas [get]
func (h *HandlerAvaliacao) BuscarEstatisticas(w http.ResponseWriter, r *http.Request) {
	faxineiroID, err := lerParamInteiro(r, "faxineiro_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "faxineiro_id inválido"))
		return
	}
	stats, err := h.servico.BuscarEstatisticasFaxineiro(r.Context(), faxineiroID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeAgregadoAvaliacao(stats))
}
