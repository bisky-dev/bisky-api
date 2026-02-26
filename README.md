# devops-dashboard API

Go API using:
- Gin (HTTP)
- pgx (Postgres driver/pool)
- sqlc-style query layer (generated code can replace the hand-written `internal/db/sqlc`)
- golang-migrate (migrations)

## Local dev (example)

1. Set env for your external Postgres:

```sh
cp api/.env.example api/.env
# then edit DATABASE_URL to your external DB
```

2. Run migrations (requires `migrate` CLI installed):

```sh
make -C api migrate-up
```

3. Run the API:

```sh
make -C api run
```

4. Hot reload (nodemon-style)

Install `air` once:

```sh
go install github.com/air-verse/air@latest
```

Then run:

```sh
make -C api dev
```

## Optional local Postgres (Docker)

If you do not want to use an external DB during development:

```sh
make postgres-up
cp api/.env.example api/.env
make -C api migrate-up
make -C api run
```

## Endpoints

- `GET /health`
