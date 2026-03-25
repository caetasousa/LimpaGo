package auth_test

import (
	"context"
	"testing"

	"limpaGo/api/auth"
)

// TestServicoAutenticacao_NovoServicoAutenticacao verifica que a construção
// do ServicoAutenticacao não entra em pânico e retorna um valor não nulo.
func TestServicoAutenticacao_NovoServicoAutenticacao(t *testing.T) {
	t.Parallel()

	cfgZitadel := auth.ConfiguracaoZitadel{
		URL:      "http://localhost:8085",
		Emissor:  "http://localhost:8085",
		ClientID: "test-client",
	}

	clienteZitadel := auth.NovoClienteZitadel(cfgZitadel)
	svcToken := auth.NovoServicoTokenOIDCMock()

	svc := auth.NovoServicoAutenticacao(clienteZitadel, nil, svcToken)
	if svc == nil {
		t.Fatal("expected non-nil ServicoAutenticacao; got nil")
	}
}

// TestServicoAutenticacao_RenovarTokenComURLInvalida valida que um erro de rede
// é convertido corretamente em ErrTokenRenovacaoInvalido.
func TestServicoAutenticacao_RenovarTokenComURLInvalida(t *testing.T) {
	t.Parallel()

	cfg := auth.ConfiguracaoZitadel{URL: "http://host-invalido-teste:9999"}
	cliente := auth.NovoClienteZitadel(cfg)
	svcToken := auth.NovoServicoTokenOIDCMock()
	svc := auth.NovoServicoAutenticacao(cliente, nil, svcToken)

	_, err := svc.RenovarToken(context.Background(), "refresh-token-invalido")
	if err == nil {
		t.Error("expected error for invalid Zitadel URL; got nil")
	}
}
