# bisky-api

Go API scaffold using:
- Gin (HTTP)
- pgx (Postgres pool/driver)
- sqlc query layer
- golang-migrate (SQL migrations)

## Requirements

- Go (1.23+)
- PostgreSQL (required)
- `migrate` CLI

Install migrate (macOS/Homebrew):

```sh
brew install golang-migrate
```

## Local Setup

From this directory (`api/`):

```sh
cp .env.example .env
```

Edit `.env` and set:
- `DATABASE_URL` (required)
- `TOKEN_ENCRYPTION_KEY` (or `PAT_ENCRYPTION_KEY`)
- `PORT` (optional, defaults to `8080`)

## Run Migrations

```sh
make migrate-up
```

Current baseline migration:
- `000001_create_users` (`users` table only)

## Run API

```sh
make run
```

`make run` auto-loads variables from `.env`.

Hot reload:

```sh
go install github.com/air-verse/air@latest
make dev
```

## Swagger / OpenAPI

Swagger docs are generated automatically by `make run` and `make dev` using `swaggo/swag`.

Install swag CLI once:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

```sh
make swagger
```

Or run API directly (auto-generates swagger first):

```sh
make run
```

Then open:

- `http://localhost:3000/swagger/index.html`
- OpenAPI JSON: `http://localhost:3000/swagger/doc.json`

## Endpoint

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/forgot-password`

Detailed endpoint docs:
- `docs/api.md`
