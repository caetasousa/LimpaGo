.PHONY: build test test-unit test-integration test-integration-zitadel test-integration-all \
        lint docker-up docker-down docker-test-up docker-test-down \
        zitadel-up zitadel-down zitadel-test-up zitadel-test-down \
        migrate swagger run run-pg dev

# Subir tudo: banco, Zitadel, migrações e API (comando principal)
dev:
	docker compose up -d postgres zitadel_postgres zitadel flyway
	@echo "Aguardando migrações..."
	@until docker inspect --format='{{.State.Status}}' limpago_flyway 2>/dev/null | grep -q exited; do sleep 1; done
	@echo "Migrações concluídas. Iniciando API..."
	DATABASE_URL="postgres://limpago:limpago_dev@localhost:5434/limpago?sslmode=disable" \
	ZITADEL_URL="http://localhost:8085" \
	ZITADEL_EMISSOR="http://localhost:8085" \
		go run ./backend/cmd/api/

# Compilar binário
build:
	go build -o limpago ./backend/cmd/api/

# Testes unitários (sem banco)
test: test-unit

test-unit:
	cd backend && go test ./... -race -count=1

# Testes de integração PostgreSQL (requer Docker)
test-integration: docker-test-up
	DATABASE_URL_TESTE="postgres://limpago:limpago_dev@localhost:5433/limpago_teste?sslmode=disable" \
		cd backend && go test ./infra/postgres/... -tags integration -race -count=1 -v
	$(MAKE) docker-test-down

# Testes de integração com Zitadel (requer Docker)
test-integration-zitadel: docker-test-up zitadel-test-up
	ZITADEL_URL_TESTE="http://localhost:8086" \
	ZITADEL_SERVICE_USER_TOKEN_TESTE="$(ZITADEL_SERVICE_USER_TOKEN_TESTE)" \
	DATABASE_URL_TESTE="postgres://limpago:limpago_dev@localhost:5433/limpago_teste?sslmode=disable" \
		cd backend && go test ./infra/zitadel/... -tags integration -race -count=1 -v
	$(MAKE) docker-test-down zitadel-test-down

# Todos os testes de integração
test-integration-all: test-integration test-integration-zitadel

# Análise estática
lint:
	cd backend && go vet ./...

# Gerar documentação Swagger
swagger:
	cd backend && swag init -g cmd/api/main.go

# Subir infraestrutura dev (postgres + zitadel + flyway)
docker-up:
	docker compose up -d postgres zitadel_postgres zitadel flyway
	docker compose logs -f flyway

# Derrubar infraestrutura dev
docker-down:
	docker compose down

# Subir apenas banco de teste (porta 5433)
docker-test-up:
	docker compose --profile test up -d postgres_teste
	@echo "Aguardando banco de teste ficar pronto..."
	@until docker exec limpago_postgres_teste pg_isready -U limpago -d limpago_teste > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Banco de teste pronto."

# Derrubar banco de teste
docker-test-down:
	docker compose --profile test down

# Subir Zitadel de produção/dev
zitadel-up:
	docker compose up -d zitadel_postgres zitadel
	@echo "Aguardando Zitadel ficar pronto..."
	@n=0; until curl -sf http://localhost:8085/debug/healthz > /dev/null 2>&1; do \
		n=$$((n+1)); \
		if [ $$n -ge 30 ]; then \
			echo "ERRO: Zitadel não ficou pronto após 60s. Logs:"; \
			docker compose logs zitadel | tail -20; \
			exit 1; \
		fi; \
		sleep 2; \
	done
	@echo "Zitadel pronto."

# Derrubar Zitadel de produção/dev
zitadel-down:
	docker compose stop zitadel zitadel_postgres

# Subir Zitadel de teste (porta 8086)
zitadel-test-up:
	docker compose --profile test up -d zitadel_teste
	@echo "Aguardando Zitadel de teste ficar pronto..."
	@n=0; until curl -sf http://localhost:8086/debug/healthz > /dev/null 2>&1; do \
		n=$$((n+1)); \
		if [ $$n -ge 30 ]; then \
			echo "ERRO: Zitadel de teste não ficou pronto após 60s. Logs:"; \
			docker compose --profile test logs zitadel_teste | tail -20; \
			exit 1; \
		fi; \
		sleep 2; \
	done
	@echo "Zitadel de teste pronto."

# Derrubar Zitadel de teste
zitadel-test-down:
	docker compose --profile test stop zitadel_teste

# Rodar API em modo desenvolvimento (in-memory, sem Zitadel)
run:
	go run ./backend/cmd/api/

# Rodar API com PostgreSQL local (sem Zitadel)
run-pg:
	DATABASE_URL="postgres://limpago:limpago_dev@localhost:5434/limpago?sslmode=disable" go run ./backend/cmd/api/

# Rodar API com PostgreSQL e Zitadel
run-pg-zitadel:
	DATABASE_URL="postgres://limpago:limpago_dev@localhost:5434/limpago?sslmode=disable" \
	ZITADEL_URL="http://localhost:8085" \
	ZITADEL_EMISSOR="http://localhost:8085" \
	ZITADEL_CLIENT_ID="$(ZITADEL_CLIENT_ID)" \
	ZITADEL_CLIENT_SECRET="$(ZITADEL_CLIENT_SECRET)" \
	ZITADEL_SERVICE_USER_TOKEN="$(ZITADEL_SERVICE_USER_TOKEN)" \
		go run ./backend/cmd/api/
