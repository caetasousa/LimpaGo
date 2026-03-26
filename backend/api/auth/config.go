package auth

import (
	"os"
	"time"
)

// ConfiguracaoJWT contém as configurações para geração e validação de tokens JWT.
type ConfiguracaoJWT struct {
	SegredoAcesso    []byte
	SegredoRenovacao []byte
	DuracaoAcesso    time.Duration
	DuracaoRenovacao time.Duration
	Emissor          string
}

// ConfiguracaoPadrao retorna a configuração JWT com valores lidos de variáveis de ambiente.
// Para desenvolvimento, usa segredos padrão caso as variáveis não estejam definidas.
func ConfiguracaoPadrao() ConfiguracaoJWT {
	segredoAcesso := os.Getenv("JWT_SEGREDO_ACESSO")
	if segredoAcesso == "" {
		segredoAcesso = "limpaGo-dev-segredo-acesso-nao-usar-em-producao"
	}

	segredoRenovacao := os.Getenv("JWT_SEGREDO_RENOVACAO")
	if segredoRenovacao == "" {
		segredoRenovacao = "limpaGo-dev-segredo-renovacao-nao-usar-em-producao"
	}

	return ConfiguracaoJWT{
		SegredoAcesso:    []byte(segredoAcesso),
		SegredoRenovacao: []byte(segredoRenovacao),
		DuracaoAcesso:    15 * time.Minute,
		DuracaoRenovacao: 7 * 24 * time.Hour,
		Emissor:          "limpaGo",
	}
}
