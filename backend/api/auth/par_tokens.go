package auth

// ParTokens contém os tokens de acesso e renovação gerados após autenticação.
type ParTokens struct {
	TokenAcesso    string `json:"token_acesso"`
	TokenRenovacao string `json:"token_renovacao"`
	TipoToken      string `json:"tipo_token"`
	ExpiraEm       int64  `json:"expira_em"`
}
