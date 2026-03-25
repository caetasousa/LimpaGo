//go:build integration

package zitadel_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"limpaGo/api/auth"
)

// criarClienteZitadelTeste cria um ClienteZitadel apontando para o Zitadel de teste.
func criarClienteZitadelTeste(t *testing.T) *auth.ClienteZitadel {
	t.Helper()
	url := os.Getenv("ZITADEL_URL_TESTE")
	if url == "" {
		t.Skip("ZITADEL_URL_TESTE não definida — pulando teste de integração Zitadel")
	}
	cfg := auth.ConfiguracaoZitadel{
		URL:              url,
		Emissor:          url,
		ClientID:         os.Getenv("ZITADEL_CLIENT_ID_TESTE"),
		ClientSecret:     os.Getenv("ZITADEL_CLIENT_SECRET_TESTE"),
		ServiceUserToken: os.Getenv("ZITADEL_SERVICE_USER_TOKEN_TESTE"),
	}
	return auth.NovoClienteZitadel(cfg)
}

// criarServicoTokenOIDCTeste cria um ServicoTokenOIDC apontando para o Zitadel de teste.
func criarServicoTokenOIDCTeste(t *testing.T) *auth.ServicoTokenOIDC {
	t.Helper()
	url := os.Getenv("ZITADEL_URL_TESTE")
	if url == "" {
		t.Skip("ZITADEL_URL_TESTE não definida — pulando teste de integração Zitadel")
	}
	cfg := auth.ConfiguracaoZitadel{
		URL:      url,
		Emissor:  url,
		ClientID: os.Getenv("ZITADEL_CLIENT_ID_TESTE"),
	}
	svc, err := auth.NovoServicoTokenOIDC(cfg)
	if err != nil {
		t.Fatalf("erro ao criar ServicoTokenOIDC de teste: %v", err)
	}
	return svc
}

// criarBancoTeste conecta ao banco de dados de teste via DATABASE_URL_TESTE.
func criarBancoTeste(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL_TESTE")
	if dsn == "" {
		t.Skip("DATABASE_URL_TESTE não definida — pulando teste de integração")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("erro ao abrir banco de teste: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("erro ao conectar ao banco de teste: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// emailTeste gera um email único para cada teste para evitar conflitos.
func emailTeste(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("teste_%s@limpago.test", sanitizarNome(t.Name()))
}

func sanitizarNome(nome string) string {
	var resultado []byte
	for _, c := range []byte(nome) {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			resultado = append(resultado, c)
		} else {
			resultado = append(resultado, '_')
		}
	}
	if len(resultado) > 40 {
		resultado = resultado[:40]
	}
	return string(resultado)
}
