package handler

import (
	"net/http"

	"limpaGo/api/dto"
	"limpaGo/domain/service"
)

// HandlerFeed gerencia o endpoint de feed de atividades.
type HandlerFeed struct {
	servico *service.ServicoFeed
}

// NovoHandlerFeed cria um novo HandlerFeed.
func NovoHandlerFeed(servico *service.ServicoFeed) *HandlerFeed {
	return &HandlerFeed{servico: servico}
}

// BuscarFeed godoc
// @Summary Buscar feed de atividades paginado
// @Tags feed
// @Produce json
// @Param pagina query int false "Número da página" default(1)
// @Param tamanho query int false "Itens por página" default(20)
// @Success 200 {object} dto.RespostaPaginaFeed
// @Router /feed [get]
func (h *HandlerFeed) BuscarFeed(w http.ResponseWriter, r *http.Request) {
	pagina := lerQueryInteiro(r, "pagina", 1)
	tamanho := lerQueryInteiro(r, "tamanho", 20)

	pagina_feed, err := h.servico.BuscarFeed(r.Context(), pagina, tamanho)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePaginaFeed(pagina_feed))
}
