.PHONY: build test test-unit test-integration lint docker-up docker-down docker-test-up docker-test-down migrate swagger run run-pg dev

# Subir banco, migrações e API (comando principal de desenvolvimento)
dev:
	docker compose up -d postgres flyway
	@echo "Aguardando migrações..."
	@until docker inspect --format='{{.State.Status}}' limpago_flyway 2>/dev/null | grep -q exited; do sleep 1; done
	@echo "Migrações concluídas. Iniciando API..."
	DATABASE_URL="postgres://limpago:limpago_dev@localhost:5434/limpago?sslmode=disable" \
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
	cd backend && DATABASE_URL_TESTE="postgres://limpago:limpago_dev@localhost:5433/limpago_teste?sslmode=disable" \
		go test ./infra/postgres/... -tags integration -race -count=1 -v
	$(MAKE) docker-test-down

# Análise estática
lint:
	cd backend && go vet ./...

# Gerar documentação Swagger
swagger:
	cd backend && swag init -g cmd/api/main.go

# Subir infraestrutura dev (postgres + flyway)
docker-up:
	docker compose up -d postgres flyway
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

# Rodar API em modo desenvolvimento (in-memory, sem banco)
run:
	go run ./backend/cmd/api/

# Rodar API com PostgreSQL local
run-pg:
	DATABASE_URL="postgres://limpago:limpago_dev@localhost:5434/limpago?sslmode=disable" go run ./backend/cmd/api/
