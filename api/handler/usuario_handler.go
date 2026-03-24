package handler

import (
	"net/http"

	"limpaGo/api/dto"
	"limpaGo/api/middleware"
	"limpaGo/domain/service"
	"limpaGo/domain/valueobject"
)

// HandlerUsuario gerencia os endpoints de usuário e perfis.
type HandlerUsuario struct {
	servico *service.ServicoUsuario
}

// NovoHandlerUsuario cria um novo HandlerUsuario.
func NovoHandlerUsuario(servico *service.ServicoUsuario) *HandlerUsuario {
	return &HandlerUsuario{servico: servico}
}

// BuscarMeuPerfil godoc
// @Summary Buscar perfil do usuário autenticado
// @Tags usuarios
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.RespostaPerfil
// @Failure 401 {object} dto.RespostaErro
// @Failure 404 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil [get]
func (h *HandlerUsuario) BuscarMeuPerfil(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	perfil, err := h.servico.BuscarPerfil(r.Context(), usuarioID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfil(perfil))
}

// AtualizarMeuPerfil godoc
// @Summary Atualizar perfil base do usuário autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoAtualizarPerfil true "Dados do perfil"
// @Success 200 {object} dto.RespostaPerfil
// @Failure 401 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil [put]
func (h *HandlerUsuario) AtualizarMeuPerfil(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoAtualizarPerfil
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	perfil, err := h.servico.AtualizarPerfil(r.Context(), usuarioID, req.NomeCompleto, req.Telefone, req.Imagem)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfil(perfil))
}

// CriarPerfilFaxineiro godoc
// @Summary Criar perfil de faxineiro para o usuário autenticado
// @Tags usuarios
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.RespostaPerfilFaxineiro
// @Failure 409 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil-faxineiro [post]
func (h *HandlerUsuario) CriarPerfilFaxineiro(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	perfil, err := h.servico.CriarPerfilFaxineiro(r.Context(), usuarioID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DePerfilFaxineiro(perfil))
}

// BuscarPerfilFaxineiro godoc
// @Summary Buscar perfil de faxineiro do usuário autenticado
// @Tags usuarios
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.RespostaPerfilFaxineiro
// @Failure 404 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil-faxineiro [get]
func (h *HandlerUsuario) BuscarPerfilFaxineiro(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	perfil, err := h.servico.BuscarPerfilFaxineiro(r.Context(), usuarioID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfilFaxineiro(perfil))
}

// AtualizarPerfilFaxineiro godoc
// @Summary Atualizar perfil de faxineiro do usuário autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoAtualizarPerfilFaxineiro true "Dados do perfil"
// @Success 200 {object} dto.RespostaPerfilFaxineiro
// @Router /usuarios/eu/perfil-faxineiro [put]
func (h *HandlerUsuario) AtualizarPerfilFaxineiro(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoAtualizarPerfilFaxineiro
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	perfil, err := h.servico.AtualizarPerfilFaxineiro(r.Context(), usuarioID, req.Descricao, req.AnosExperiencia, req.Especialidades, req.CidadesAtendidas)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfilFaxineiro(perfil))
}

// CriarPerfilCliente godoc
// @Summary Criar perfil de cliente para o usuário autenticado
// @Tags usuarios
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.RespostaPerfilCliente
// @Failure 409 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil-cliente [post]
func (h *HandlerUsuario) CriarPerfilCliente(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	perfil, err := h.servico.CriarPerfilCliente(r.Context(), usuarioID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusCreated, dto.DePerfilCliente(perfil))
}

// BuscarPerfilCliente godoc
// @Summary Buscar perfil de cliente do usuário autenticado
// @Tags usuarios
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.RespostaPerfilCliente
// @Failure 404 {object} dto.RespostaErro
// @Router /usuarios/eu/perfil-cliente [get]
func (h *HandlerUsuario) BuscarPerfilCliente(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	perfil, err := h.servico.BuscarPerfilCliente(r.Context(), usuarioID)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfilCliente(perfil))
}

// AtualizarPerfilCliente godoc
// @Summary Atualizar perfil de cliente do usuário autenticado
// @Tags usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RequisicaoAtualizarPerfilCliente true "Dados do perfil"
// @Success 200 {object} dto.RespostaPerfilCliente
// @Router /usuarios/eu/perfil-cliente [put]
func (h *HandlerUsuario) AtualizarPerfilCliente(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := middleware.ObterUsuarioID(r.Context())
	if !ok {
		escreverJSON(w, http.StatusUnauthorized, dto.NovaRespostaErro(http.StatusUnauthorized, "não autenticado"))
		return
	}
	var req dto.RequisicaoAtualizarPerfilCliente
	if err := lerJSON(r, &req); err != nil {
		escreverJSON(w, http.StatusBadRequest, dto.NovaRespostaErro(http.StatusBadRequest, "corpo inválido"))
		return
	}
	perfil, err := h.servico.AtualizarPerfilCliente(r.Context(), usuarioID,
		req.Endereco.ParaEndereco(),
		valueobject.TipoImovel(req.TipoImovel),
		req.Quartos, req.Banheiros,
		req.TamanhoImovelM2, req.Observacoes,
	)
	if err != nil {
		escreverErro(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, dto.DePerfilCliente(perfil))
}
