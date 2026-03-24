package handler

import (
	"net/http"
	"time"

	"limpaGo/api/dto"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
)

// HandlerAgenda gerencia os endpoints de disponibilidade e bloqueios de agenda.
type HandlerAgenda struct {
	servico *service.ServicoAgenda
}

// NovoHandlerAgenda cria um novo HandlerAgenda.
func NovoHandlerAgenda(servico *service.ServicoAgenda) *HandlerAgenda {
	return &HandlerAgenda{servico: servico}
}

// ListarDisponibilidade godoc
// @Summary Listar disponibilidades do faxineiro autenticado
// @Tags agenda
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RespostaDisponibilidade
// @Failure 401 {object} dto.RespostaErro
// @Router /agenda/disponibilidades [get]
func (h *HandlerAgenda) ListarDisponibilidade(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	lista, err := h.servico.ListarDisponibilidade(r.Context(), faxineiroID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeDisponibilidadeLista(lista))
}

// AdicionarDisponibilidade godoc
// @Summary Adicionar bloco de disponibilidade
// @Tags agenda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoDisponibilidade true "Dados da disponibilidade"
// @Success 201 {object} dto.RespostaDisponibilidade
// @Failure 401 {object} dto.RespostaErro
// @Failure 422 {object} dto.RespostaErro
// @Router /agenda/disponibilidades [post]
func (h *HandlerAgenda) AdicionarDisponibilidade(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoDisponibilidade
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	d, err := h.servico.AdicionarDisponibilidade(r.Context(), faxineiroID, time.Weekday(req.DiaSemana), req.HoraInicio, req.HoraFim)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DeDisponibilidade(d))
}

// RemoverDisponibilidade godoc
// @Summary Remover bloco de disponibilidade
// @Tags agenda
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da disponibilidade"
// @Success 204
// @Failure 401 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /agenda/disponibilidades/{id} [delete]
func (h *HandlerAgenda) RemoverDisponibilidade(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	id, err := lerParamInteiro(r, "id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "id inválido"))
		return
	}
	if err := h.servico.RemoverDisponibilidade(r.Context(), id, faxineiroID); err != nil {
		escreverErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListarBloqueios godoc
// @Summary Listar bloqueios do faxineiro autenticado
// @Tags agenda
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RespostaBloqueio
// @Failure 401 {object} dto.RespostaErro
// @Router /agenda/bloqueios [get]
func (h *HandlerAgenda) ListarBloqueios(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	lista, err := h.servico.ListarBloqueios(r.Context(), faxineiroID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeBloqueioLista(lista))
}

// CriarBloqueioPessoal godoc
// @Summary Criar bloqueio pessoal de agenda
// @Tags agenda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoBloqueio true "Dados do bloqueio"
// @Success 201 {object} dto.RespostaBloqueio
// @Failure 401 {object} dto.RespostaErro
// @Failure 422 {object} dto.RespostaErro
// @Router /agenda/bloqueios [post]
func (h *HandlerAgenda) CriarBloqueioPessoal(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoBloqueio
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	bloqueio, err := h.servico.CriarBloqueioPessoal(r.Context(), faxineiroID, req.DataInicio, req.DataFim)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DeBloqueio(bloqueio))
}

// RemoverBloqueioPessoal godoc
// @Summary Remover bloqueio pessoal de agenda
// @Tags agenda
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do bloqueio"
// @Success 204
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /agenda/bloqueios/{id} [delete]
func (h *HandlerAgenda) RemoverBloqueioPessoal(w http.ResponseWriter, r *http.Request) {
	faxineiroID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	id, err := lerParamInteiro(r, "id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "id inválido"))
		return
	}
	if err := h.servico.RemoverBloqueioPessoal(r.Context(), id, faxineiroID); err != nil {
		escreverErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
