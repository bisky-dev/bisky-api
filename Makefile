.PHONY: run dev lint lint-types lint-go sqlc swagger migrate-up migrate-down postgres-up postgres-down

ENV_LOADER=if [ -f .env ]; then set -a; . ./.env; set +a; fi;
SWAG_BIN=$(shell if command -v swag >/dev/null 2>&1; then command -v swag; else echo "$$(go env GOPATH)/bin/swag"; fi)
AIR_BIN=$(shell if command -v air >/dev/null 2>&1; then command -v air; else echo "$$(go env GOPATH)/bin/air"; fi)

run:
	@$(MAKE) swagger
	@echo "Starting API on :$${PORT:-8080}"
	@$(ENV_LOADER) PORT=$${PORT:-$${PORT:-8080}} go run ./cmd/api

dev:
	@$(MAKE) swagger
	@echo "Starting API with hot reload on :$${PORT:-8080}"
	@if [ ! -x "$(AIR_BIN)" ]; then echo "air not found. Install with: go install github.com/air-verse/air@latest"; exit 1; fi
	@$(ENV_LOADER) PORT=$${PORT:-$${PORT:-8080}} $(AIR_BIN) -c .air.toml

lint: lint-types lint-go

lint-types:
	go run ./tools/typefilelint

lint-go:
	golangci-lint run

sqlc:
	sqlc generate

swagger:
	@if [ ! -x "$(SWAG_BIN)" ]; then echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@v1.16.6"; exit 1; fi
	$(SWAG_BIN) init -g cmd/api/main.go -o docs/swagger --parseInternal

migrate-up:
	@$(ENV_LOADER) migrate -path db/migrations -database "$$DATABASE_URL" up

migrate-down:
	@$(ENV_LOADER) migrate -path db/migrations -database "$$DATABASE_URL" down 1

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose stop postgres
