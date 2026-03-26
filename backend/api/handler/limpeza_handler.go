package handler

import (
	"net/http"

	"limpaGo/api/dto"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
	"limpaGo/domain/valueobject"
)

// HandlerLimpeza gerencia os endpoints de serviços de limpeza.
type HandlerLimpeza struct {
	servico *service.ServicoLimpeza
}

// NovoHandlerLimpeza cria um novo HandlerLimpeza.
func NovoHandlerLimpeza(servico *service.ServicoLimpeza) *HandlerLimpeza {
	return &HandlerLimpeza{servico: servico}
}

// ListarCatalogo godoc
// @Summary Listar catálogo de limpezas
// @Tags limpezas
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param tamanho query int false "Itens por página" default(20)
// @Success 200 {array} dto.RespostaLimpeza
// @Router /limpezas [get]
func (h *HandlerLimpeza) ListarCatalogo(w http.ResponseWriter, r *http.Request) {
	pagina := lerQueryInteiro(r, "pagina", 1)
	tamanho := lerQueryInteiro(r, "tamanho", 20)

	limpezas, err := h.servico.ListarCatalogo(r.Context(), pagina, tamanho)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeLimpezaLista(limpezas))
}

// BuscarLimpeza godoc
// @Summary Buscar limpeza por ID
// @Tags limpezas
// @Produce json
// @Param id path int true "ID da limpeza"
// @Success 200 {object} dto.RespostaLimpeza
// @Failure 404 {object} dto.RespostaErro
// @Router /limpezas/{id} [get]
func (h *HandlerLimpeza) BuscarLimpeza(w http.ResponseWriter, r *http.Request) {
	id, err := lerParamInteiro(r, "id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "id inválido"))
		return
	}
	limpeza, err := h.servico.BuscarPorID(r.Context(), id)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeLimpeza(limpeza))
}

// CriarLimpeza godoc
// @Summary Criar novo serviço de limpeza
// @Tags limpezas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoCriarLimpeza true "Dados da limpeza"
// @Success 201 {object} dto.RespostaLimpeza
// @Failure 401 {object} dto.RespostaErro
// @Failure 422 {object} dto.RespostaErro
// @Router /limpezas [post]
func (h *HandlerLimpeza) CriarLimpeza(w http.ResponseWriter, r *http.Request) {
	profissionalID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoCriarLimpeza
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	limpeza, err := h.servico.Criar(r.Context(), profissionalID, req.Nome, req.Descricao, req.ValorHora, req.DuracaoEstimada, valueobject.TipoLimpeza(req.TipoLimpeza))
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DeLimpeza(limpeza))
}

// AtualizarLimpeza godoc
// @Summary Atualizar serviço de limpeza
// @Tags limpezas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da limpeza"
// @Param body body dto.RequisicaoAtualizarLimpeza true "Dados da limpeza"
// @Success 200 {object} dto.RespostaLimpeza
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /limpezas/{id} [put]
func (h *HandlerLimpeza) AtualizarLimpeza(w http.ResponseWriter, r *http.Request) {
	profissionalID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	id, err := lerParamInteiro(r, "id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "id inválido"))
		return
	}
	var req dto.RequisicaoAtualizarLimpeza
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	limpeza, err := h.servico.Atualizar(r.Context(), id, profissionalID, req.Nome, req.Descricao, req.ValorHora, req.DuracaoEstimada, valueobject.TipoLimpeza(req.TipoLimpeza))
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeLimpeza(limpeza))
}

// DeletarLimpeza godoc
// @Summary Deletar serviço de limpeza
// @Tags limpezas
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID da limpeza"
// @Success 204
// @Failure 401 {object} dto.RespostaErro
// @Failure 403 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /limpezas/{id} [delete]
func (h *HandlerLimpeza) DeletarLimpeza(w http.ResponseWriter, r *http.Request) {
	profissionalID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	id, err := lerParamInteiro(r, "id")
	if err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "id inválido"))
		return
	}
	if err := h.servico.Deletar(r.Context(), id, profissionalID); err != nil {
		escreverErro(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListarMinhasLimpezas godoc
// @Summary Listar limpezas do profissional autenticado
// @Tags limpezas
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RespostaLimpeza
// @Failure 401 {object} dto.RespostaErro
// @Router /usuarios/eu/limpezas [get]
func (h *HandlerLimpeza) ListarMinhasLimpezas(w http.ResponseWriter, r *http.Request) {
	profissionalID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	limpezas, err := h.servico.ListarPorProfissional(r.Context(), profissionalID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DeLimpezaLista(limpezas))
}
