//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// criarBancoTeste conecta ao banco de teste e executa todas as migrações.
// Requer que DATABASE_URL_TESTE esteja definida (ex: via Makefile).
func criarBancoTeste(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL_TESTE")
	if dsn == "" {
		t.Skip("DATABASE_URL_TESTE não definida — pulando teste de integração")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("abrir banco de teste: %v", err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("conectar ao banco de teste: %v", err)
	}

	executarMigracoes(t, db)

	t.Cleanup(func() { db.Close() })

	return db
}

// executarMigracoes lê e executa todos os arquivos SQL em db/migrations/ em ordem.
func executarMigracoes(t *testing.T, db *sql.DB) {
	t.Helper()

	_, arquivo, _, _ := runtime.Caller(0)
	raiz := filepath.Join(filepath.Dir(arquivo), "..", "..", "db", "migrations")

	entradas, err := os.ReadDir(raiz)
	if err != nil {
		t.Fatalf("ler diretório de migrações %s: %v", raiz, err)
	}

	var arquivosSql []string
	for _, e := range entradas {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			arquivosSql = append(arquivosSql, filepath.Join(raiz, e.Name()))
		}
	}
	// Ordenar pelo número da versão (V1, V2, ..., V10) — não alfabeticamente
	sort.Slice(arquivosSql, func(i, j int) bool {
		return versaoMigracao(arquivosSql[i]) < versaoMigracao(arquivosSql[j])
	})

	for _, caminho := range arquivosSql {
		conteudo, err := os.ReadFile(caminho)
		if err != nil {
			t.Fatalf("ler migração %s: %v", caminho, err)
		}
		if _, err := db.Exec(string(conteudo)); err != nil {
			// Ignorar erros de "tabela já existe" para permitir reutilização do banco
			if !strings.Contains(err.Error(), "already exists") {
				t.Fatalf("executar migração %s: %v", filepath.Base(caminho), err)
			}
		}
	}
}

// limparTabelas remove todos os dados das tabelas em ordem segura (respeitando FK).
func limparTabelas(t *testing.T, db *sql.DB) {
	t.Helper()

	tabelas := []string{
		"avaliacoes",
		"bloqueios",
		"disponibilidades",
		"solicitacoes",
		"limpezas",
		"perfis_cliente",
		"perfis_profissional",
		"perfis",
		"credenciais",
		"usuarios",
	}

	for _, tabela := range tabelas {
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tabela)); err != nil {
			t.Fatalf("limpar tabela %s: %v", tabela, err)
		}
	}
}

// inserirUsuario insere um usuário diretamente no banco e retorna o ID gerado.
func inserirUsuario(t *testing.T, db *sql.DB, email, nomeUsuario string) int {
	t.Helper()

	var id int
	err := db.QueryRow(
		`INSERT INTO usuarios (email, nome_usuario) VALUES ($1, $2) RETURNING id`,
		email, nomeUsuario,
	).Scan(&id)
	if err != nil {
		t.Fatalf("inserir usuário %s: %v", email, err)
	}
	return id
}

// inserirLimpeza insere uma limpeza diretamente no banco e retorna o ID.
func inserirLimpeza(t *testing.T, db *sql.DB, profissionalID int, nome string) int {
	t.Helper()

	var id int
	err := db.QueryRow(
		`INSERT INTO limpezas (nome, descricao, valor_hora, duracao_estimada, tipo_limpeza, profissional_id)
		 VALUES ($1, '', 50.00, 2.0, 'limpeza_padrao', $2) RETURNING id`,
		nome, profissionalID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("inserir limpeza %s: %v", nome, err)
	}
	return id
}

// versaoMigracao extrai o número inteiro da versão de um arquivo Flyway (ex: "V10__..." → 10).
func versaoMigracao(caminho string) int {
	nome := filepath.Base(caminho)
	// Remove o prefixo "V" e pega tudo até o primeiro "_"
	nome = strings.TrimPrefix(nome, "V")
	partes := strings.SplitN(nome, "_", 2)
	if len(partes) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(partes[0])
	return n
}

// inserirSolicitacao insere uma solicitação diretamente no banco e retorna o ID.
func inserirSolicitacao(t *testing.T, db *sql.DB, clienteID, limpezaID int) int {
	t.Helper()

	var id int
	err := db.QueryRow(
		`INSERT INTO solicitacoes
		 (cliente_id, limpeza_id, status, data_agendada, preco_total, multa_cancelamento,
		  endereco_rua, endereco_complemento, endereco_bairro, endereco_cidade, endereco_estado, endereco_cep)
		 VALUES ($1, $2, 'pendente', NOW() + INTERVAL '2 days', 100.00, 0,
		         'Rua Teste', '', 'Centro', 'São Paulo', 'SP', '01310-100')
		 RETURNING id`,
		clienteID, limpezaID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("inserir solicitação: %v", err)
	}
	return id
}
