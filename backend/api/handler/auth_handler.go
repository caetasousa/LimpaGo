package handler

import (
	"net/http"

	"limpaGo/api/auth"
	"limpaGo/api/dto"
)

// HandlerAutenticacao gerencia os endpoints de autenticação.
type HandlerAutenticacao struct {
	servico *auth.ServicoAutenticacao
	cfgOIDC auth.ConfiguracaoZitadel
}

// NovoHandlerAutenticacao cria um novo HandlerAutenticacao.
func NovoHandlerAutenticacao(servico *auth.ServicoAutenticacao, cfg auth.ConfiguracaoZitadel) *HandlerAutenticacao {
	return &HandlerAutenticacao{servico: servico, cfgOIDC: cfg}
}

// Registrar godoc
// @Summary      Registrar novo usuário
// @Description  Cria uma nova conta com email, nome de usuário e senha via Zitadel
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RequisicaoRegistroComSenha  true  "Dados de registro"
// @Success      201   {object}  dto.RespostaAutenticacao
// @Failure      409   {object}  dto.RespostaErro
// @Failure      422   {object}  dto.RespostaErro
// @Router       /auth/registrar [post]
func (h *HandlerAutenticacao) Registrar(w http.ResponseWriter, r *http.Request) {
	var req dto.RequisicaoRegistroComSenha
	if err := lerJSON(r, &req); err != nil {
		escreverErro(w, err)
		return
	}

	usuario, tokens, err := h.servico.Registrar(r.Context(), req.Email, req.NomeUsuario, req.Senha)
	if err != nil {
		escreverErro(w, err)
		return
	}

	resp := dto.RespostaAutenticacao{
		Usuario: dto.DeUsuario(usuario),
		Tokens:  dto.DeParTokens(tokens),
	}
	escreverJSON(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary      Login
// @Description  Autentica o usuário com email e senha via Zitadel e retorna tokens JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RequisicaoLogin  true  "Credenciais de login"
// @Success      200   {object}  dto.RespostaAutenticacao
// @Failure      401   {object}  dto.RespostaErro
// @Failure      403   {object}  dto.RespostaErro
// @Router       /auth/login [post]
func (h *HandlerAutenticacao) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.RequisicaoLogin
	if err := lerJSON(r, &req); err != nil {
		escreverErro(w, err)
		return
	}

	usuario, tokens, err := h.servico.Login(r.Context(), req.Email, req.Senha)
	if err != nil {
		escreverErro(w, err)
		return
	}

	resp := dto.RespostaAutenticacao{
		Usuario: dto.DeUsuario(usuario),
		Tokens:  dto.DeParTokens(tokens),
	}
	escreverJSON(w, http.StatusOK, resp)
}

// RenovarToken godoc
// @Summary      Renovar token
// @Description  Usa o token de renovação do Zitadel para obter um novo par de tokens JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RequisicaoRenovarToken  true  "Token de renovação"
// @Success      200   {object}  dto.ParTokensDTO
// @Failure      401   {object}  dto.RespostaErro
// @Router       /auth/renovar [post]
func (h *HandlerAutenticacao) RenovarToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RequisicaoRenovarToken
	if err := lerJSON(r, &req); err != nil {
		escreverErro(w, err)
		return
	}

	tokens, err := h.servico.RenovarToken(r.Context(), req.TokenRenovacao)
	if err != nil {
		escreverErro(w, err)
		return
	}

	escreverJSON(w, http.StatusOK, dto.DeParTokens(tokens))
}

// ConfiguracaoOIDC godoc
// @Summary      Configuração OIDC
// @Description  Retorna a URL e client_id do Zitadel para configuração do frontend
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.RespostaConfiguracaoOIDC
// @Router       /auth/config [get]
func (h *HandlerAutenticacao) ConfiguracaoOIDC(w http.ResponseWriter, r *http.Request) {
	escreverJSON(w, http.StatusOK, dto.RespostaConfiguracaoOIDC{
		URL:      h.cfgOIDC.URL,
		ClientID: h.cfgOIDC.ClientID,
		Emissor:  h.cfgOIDC.Emissor,
	})
}
