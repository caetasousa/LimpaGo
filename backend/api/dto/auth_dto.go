package dto

import "limpaGo/api/auth"

// RequisicaoRegistroComSenha representa o corpo da requisição de registro com senha.
type RequisicaoRegistroComSenha struct {
	Email       string `json:"email"`
	NomeUsuario string `json:"nome_usuario"`
	Senha       string `json:"senha"`
}

// RequisicaoLogin representa o corpo da requisição de login.
type RequisicaoLogin struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// RequisicaoRenovarToken representa o corpo da requisição de renovação de token.
type RequisicaoRenovarToken struct {
	TokenRenovacao string `json:"token_renovacao"`
}

// ParTokensDTO representa o par de tokens JWT retornado ao cliente.
type ParTokensDTO struct {
	TokenAcesso    string `json:"token_acesso"`
	TokenRenovacao string `json:"token_renovacao"`
	TipoToken      string `json:"tipo_token"`
	ExpiraEm       int64  `json:"expira_em"`
}

// RespostaAutenticacao é retornada após registro ou login bem-sucedido.
type RespostaAutenticacao struct {
	Usuario RespostaUsuario `json:"usuario"`
	Tokens  ParTokensDTO    `json:"tokens"`
}

// RespostaConfiguracaoOIDC contém os dados de configuração do provedor OIDC para o frontend.
type RespostaConfiguracaoOIDC struct {
	URL      string `json:"url"`
	ClientID string `json:"client_id"`
	Emissor  string `json:"emissor"`
}

// DeParTokens converte um ParTokens do domínio auth em ParTokensDTO.
func DeParTokens(p *auth.ParTokens) ParTokensDTO {
	return ParTokensDTO{
		TokenAcesso:    p.TokenAcesso,
		TokenRenovacao: p.TokenRenovacao,
		TipoToken:      p.TipoToken,
		ExpiraEm:       p.ExpiraEm,
	}
}
