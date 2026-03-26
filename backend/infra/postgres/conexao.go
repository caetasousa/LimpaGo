// Package postgres contém as implementações de repositório para PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ConfiguracaoBanco agrupa os parâmetros de conexão com o banco de dados.
type ConfiguracaoBanco struct {
	URL             string
	Host            string
	Port            int
	Usuario         string
	Senha           string
	Banco           string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// CarregarConfiguracaoBanco lê a configuração do banco a partir de variáveis de ambiente.
// DATABASE_URL tem precedência sobre as variáveis individuais.
func CarregarConfiguracaoBanco() ConfiguracaoBanco {
	cfg := ConfiguracaoBanco{
		URL:             os.Getenv("DATABASE_URL"),
		Host:            getEnvOuPadrao("PG_HOST", "localhost"),
		Port:            getEnvInteiroOuPadrao("PG_PORT", 5432),
		Usuario:         getEnvOuPadrao("PG_USER", "limpago"),
		Senha:           getEnvOuPadrao("PG_PASSWORD", "limpago_dev"),
		Banco:           getEnvOuPadrao("PG_DATABASE", "limpago"),
		SSLMode:         getEnvOuPadrao("PG_SSLMODE", "disable"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
	return cfg
}

// DSN retorna a string de conexão.
func (c ConfiguracaoBanco) DSN() string {
	if c.URL != "" {
		return c.URL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Usuario, c.Senha, c.Banco, c.SSLMode,
	)
}

// NovoBanco abre a conexão com PostgreSQL e verifica a conectividade.
func NovoBanco(cfg ConfiguracaoBanco) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("abrir conexão com banco: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("verificar conexão com banco: %w", err)
	}

	return db, nil
}

func getEnvOuPadrao(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

func getEnvInteiroOuPadrao(chave string, padrao int) int {
	if v := os.Getenv(chave); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return padrao
}
