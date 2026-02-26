.PHONY: run dev lint lint-types lint-go sqlc swagger migrate-up migrate-down postgres-up postgres-down

ENV_LOADER=if [ -f .env ]; then set -a; . ./.env; set +a; fi;

run:
	@$(MAKE) swagger
	@echo "Starting API on :$${PORT:-8080}"
	@$(ENV_LOADER) PORT=$${PORT:-$${PORT:-8080}} go run ./cmd/api

dev:
	@$(MAKE) swagger
	@echo "Starting API with hot reload on :$${PORT:-8080}"
	@$(ENV_LOADER) PORT=$${PORT:-$${PORT:-8080}} air -c .air.toml

lint: lint-types lint-go

lint-types:
	go run ./tools/typefilelint

lint-go:
	golangci-lint run

sqlc:
	sqlc generate

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs/swagger --parseInternal

migrate-up:
	@$(ENV_LOADER) migrate -path db/migrations -database "$$DATABASE_URL" up

migrate-down:
	@$(ENV_LOADER) migrate -path db/migrations -database "$$DATABASE_URL" down 1

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose stop postgres
