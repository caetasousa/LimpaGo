package auth

import "os"

// ConfiguracaoZitadel contém as configurações para integração com o Zitadel OIDC.
type ConfiguracaoZitadel struct {
	URL               string // URL base do Zitadel, ex: http://localhost:8085
	Emissor           string // Issuer esperado nos tokens JWT
	ClientID          string // ID da aplicação registrada no Zitadel
	ClientSecret      string // Secret da aplicação (usado no fluxo OAuth2)
	ServiceUserToken  string // PAT do service user para a Management API
}

// CarregarConfiguracaoZitadel retorna a configuração Zitadel lida das variáveis de ambiente.
func CarregarConfiguracaoZitadel() ConfiguracaoZitadel {
	url := os.Getenv("ZITADEL_URL")
	if url == "" {
		url = "http://localhost:8085"
	}

	emissor := os.Getenv("ZITADEL_EMISSOR")
	if emissor == "" {
		emissor = url
	}

	return ConfiguracaoZitadel{
		URL:              url,
		Emissor:          emissor,
		ClientID:         os.Getenv("ZITADEL_CLIENT_ID"),
		ClientSecret:     os.Getenv("ZITADEL_CLIENT_SECRET"),
		ServiceUserToken: os.Getenv("ZITADEL_SERVICE_USER_TOKEN"),
	}
}
