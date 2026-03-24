package handler

import (
	"net/http"

	"limpaGo/api/dto"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
)

// HandlerSolicitacao gerencia os endpoints de solicitações de limpeza.
type HandlerSolicitacao struct {
	servico *service.ServicoSolicitacao
}

// NovoHandlerSolicitacao cria um novo HandlerSolicitacao.
func NovoHandlerSolicitacao(servico *service.ServicoSolicitacao) *HandlerSolicitacao {
	return &HandlerSolicitacao{servico: servico}
}

// CriarSolicitacao godoc
// @Summary Criar solicitação de limpeza
// @Tags solicitacoes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoCriarSolicitacao true "Dados da solicitação"
// @Success 201 {object} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 409 {object} dto.RespostaErro
// @Failure 422 {object} dto.RespostaErro
// @Router /solicitacoes [post]
func (h *HandlerSolicitacao) CriarSolicitacao(w http.ResponseWriter, r *http.Request) {
	clienteID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoCriarSolicitacao
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	solicitacao, err := h.servico.CriarSolicitacao(r.Context(), clienteID, req.LimpezaID, req.DataAgendada)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DeSolicitacao(solicitacao))
}

// AceitarSolicitacao godoc
// @Summary Aceitar solicitação de limpeza
// @Tags solicitacoes
// @Produce json
// @Security BearerAuth
// @Param cliente_id path int true "ID do cliente"
// @Param limpeza_id path int true "ID da limpeza"
// @Success 200 {object} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /solicitacoes/{cliente_id}/{limpeza_id}/aceitar [post]
func (h *HandlerSolicitacao) AceitarSolicitacao(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	clienteID, err := lerParamInteiro(r, "cliente_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "cliente_id inválido"))
		return
	}
	limpezaID, err := lerParamInteiro(r, "limpeza_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "limpeza_id inválido"))
		return
	}
	solicitacao, err := h.servico.AceitarSolicitacao(r.Context(), faxineiroID, clienteID, limpezaID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeSolicitacao(solicitacao))
}

// RejeitarSolicitacao godoc
// @Summary Rejeitar solicitação de limpeza
// @Tags solicitacoes
// @Produce json
// @Security BearerAuth
// @Param cliente_id path int true "ID do cliente"
// @Param limpeza_id path int true "ID da limpeza"
// @Success 200 {object} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /solicitacoes/{cliente_id}/{limpeza_id}/rejeitar [post]
func (h *HandlerSolicitacao) RejeitarSolicitacao(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	clienteID, err := lerParamInteiro(r, "cliente_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "cliente_id inválido"))
		return
	}
	limpezaID, err := lerParamInteiro(r, "limpeza_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "limpeza_id inválido"))
		return
	}
	solicitacao, err := h.servico.RejeitarSolicitacao(r.Context(), faxineiroID, clienteID, limpezaID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeSolicitacao(solicitacao))
}

// CancelarSolicitacao godoc
// @Summary Cancelar solicitação de limpeza
// @Tags solicitacoes
// @Produce json
// @Security BearerAuth
// @Param limpeza_id path int true "ID da limpeza"
// @Success 200 {object} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /solicitacoes/{limpeza_id}/cancelar [post]
func (h *HandlerSolicitacao) CancelarSolicitacao(w http.ResponseWriter, r *http.Request) {
	clienteID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	limpezaID, err := lerParamInteiro(r, "limpeza_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "limpeza_id inválido"))
		return
	}
	solicitacao, err := h.servico.CancelarSolicitacao(r.Context(), clienteID, limpezaID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeSolicitacao(solicitacao))
}

// ListarPorLimpeza godoc
// @Summary Listar solicitações de uma limpeza (faxineiro)
// @Tags solicitacoes
// @Produce json
// @Security BearerAuth
// @Param limpeza_id path int true "ID da limpeza"
// @Success 200 {array} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Router /limpezas/{limpeza_id}/solicitacoes [get]
func (h *HandlerSolicitacao) ListarPorLimpeza(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	limpezaID, err := lerParamInteiro(r, "limpeza_id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "limpeza_id inválido"))
		return
	}
	solicitacoes, err := h.servico.ListarSolicitacoesPorLimpeza(r.Context(), faxineiroID, limpezaID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeSolicitacaoLista(solicitacoes))
}

// ListarMinhasSolicitacoes godoc
// @Summary Listar solicitações do cliente autenticado
// @Tags solicitacoes
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RespostaSolicitacao
// @Failure 401 {object} dto.RespostaErro
// @Router /usuarios/eu/solicitacoes [get]
func (h *HandlerSolicitacao) ListarMinhasSolicitacoes(w http.ResponseWriter, r *http.Request) {
	clienteID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	solicitacoes, err := h.servico.ListarSolicitacoesPorCliente(r.Context(), clienteID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeSolicitacaoLista(solicitacoes))
}
