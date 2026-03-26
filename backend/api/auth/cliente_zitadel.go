package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TokensZitadel representa os tokens retornados pelo Zitadel após autenticação.
type TokensZitadel struct {
	TokenAcesso    string
	TokenRenovacao string
	ExpiraEm       int64
}

// ClienteZitadel encapsula as chamadas HTTP para a API do Zitadel.
type ClienteZitadel struct {
	urlBase          string
	clientID         string
	clientSecret     string
	serviceUserToken string
	httpClient       *http.Client
}

// NovoClienteZitadel cria um ClienteZitadel com a configuração fornecida.
func NovoClienteZitadel(cfg ConfiguracaoZitadel) *ClienteZitadel {
	return &ClienteZitadel{
		urlBase:          cfg.URL,
		clientID:         cfg.ClientID,
		clientSecret:     cfg.ClientSecret,
		serviceUserToken: cfg.ServiceUserToken,
		httpClient:       &http.Client{Timeout: 10 * time.Second},
	}
}

// RegistrarUsuario cria um novo usuário humano no Zitadel via Management API
// e retorna o ID externo (sub) do usuário criado.
func (c *ClienteZitadel) RegistrarUsuario(ctx context.Context, email, nomeUsuario, senha string) (string, error) {
	corpo := map[string]any{
		"userName": nomeUsuario,
		"profile": map[string]any{
			"displayName": nomeUsuario,
		},
		"email": map[string]any{
			"email":           email,
			"isEmailVerified": false,
		},
		"password": map[string]any{
			"password":       senha,
			"changeRequired": false,
		},
	}

	resposta, err := c.chamarAPI(ctx, http.MethodPost, "/v2/users/human", corpo)
	if err != nil {
		return "", fmt.Errorf("registrar usuário no Zitadel: %w", err)
	}

	idExterno, ok := resposta["userId"].(string)
	if !ok || idExterno == "" {
		return "", fmt.Errorf("zitadel não retornou userId: %w", ErrCredenciaisInvalidas)
	}

	return idExterno, nil
}

// Autenticar valida email e senha via Zitadel e retorna tokens OAuth2.
func (c *ClienteZitadel) Autenticar(ctx context.Context, email, senha string) (*TokensZitadel, error) {
	return c.obterTokens(ctx, url.Values{
		"grant_type": {"password"},
		"username":   {email},
		"password":   {senha},
		"scope":      {"openid profile email offline_access"},
	})
}

// RenovarToken usa um refresh_token para obter novos tokens via Zitadel.
func (c *ClienteZitadel) RenovarToken(ctx context.Context, refreshToken string) (*TokensZitadel, error) {
	return c.obterTokens(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"scope":         {"openid profile email offline_access"},
	})
}

// obterTokens realiza um Token Request OAuth2 no endpoint /oauth/v2/token do Zitadel.
func (c *ClienteZitadel) obterTokens(ctx context.Context, params url.Values) (*TokensZitadel, error) {
	endpoint := fmt.Sprintf("%s/oauth/v2/token", c.urlBase)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("criar requisição de token: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chamar token endpoint: %w", err)
	}
	defer resp.Body.Close()

	var resultado map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return nil, fmt.Errorf("decodificar resposta de token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		descricao, _ := resultado["error_description"].(string)
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusBadRequest {
			return nil, ErrCredenciaisInvalidas
		}
		return nil, fmt.Errorf("erro do Zitadel (%d): %s", resp.StatusCode, descricao)
	}

	tokenAcesso, _ := resultado["access_token"].(string)
	tokenRenovacao, _ := resultado["refresh_token"].(string)
	expiraEmSeg, _ := resultado["expires_in"].(float64)

	if tokenAcesso == "" {
		return nil, ErrCredenciaisInvalidas
	}

	return &TokensZitadel{
		TokenAcesso:    tokenAcesso,
		TokenRenovacao: tokenRenovacao,
		ExpiraEm:       time.Now().Add(time.Duration(expiraEmSeg) * time.Second).Unix(),
	}, nil
}

// chamarAPI realiza uma chamada autenticada à API de gestão do Zitadel.
func (c *ClienteZitadel) chamarAPI(ctx context.Context, metodo, caminho string, corpo any) (map[string]any, error) {
	corpoJSON, err := json.Marshal(corpo)
	if err != nil {
		return nil, fmt.Errorf("serializar corpo: %w", err)
	}

	endpoint := fmt.Sprintf("%s%s", c.urlBase, caminho)
	req, err := http.NewRequestWithContext(ctx, metodo, endpoint, bytes.NewReader(corpoJSON))
	if err != nil {
		return nil, fmt.Errorf("criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.serviceUserToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executar requisição: %w", err)
	}
	defer resp.Body.Close()

	var resultado map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return nil, fmt.Errorf("decodificar resposta: %w", err)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode == http.StatusConflict {
			return nil, ErrEmailJaCadastradoNoIdP
		}
		msg, _ := resultado["message"].(string)
		return nil, fmt.Errorf("erro do Zitadel (%d): %s", resp.StatusCode, msg)
	}

	return resultado, nil
}
