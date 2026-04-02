# Infrastructure Memory

## PostgreSQL Connection Pool

- Driver: `github.com/jackc/pgx/v5/pgxpool`
- `MaxConns: 20`
- `MinConns: 4`
- `MaxConnLifetime: 30m`
- `MaxConnIdleTime: 5m`
- Opened in `main.go`; passed to `ServerBuilder.WithPool()`
- Connection verified by a ping at startup; process exits if unreachable

## Redis

- Client: `github.com/redis/go-redis/v9`
- AOF persistence enabled in `docker-compose.yml` (`--appendonly yes`) — Redis also carries the job queue
- Opened in `main.go`; passed to `ServerBuilder.WithRedis()`
- Connection verified by a ping at startup; process exits if unreachable
- Local URL: `redis://redis:6379/0`

## Migrations

- Tool: `github.com/pressly/goose/v3`
- Migration files: `backend/migrations/*.sql` (goose timestamp format)
- Embedded in binary via `backend/migrations/embed.go` (`//go:embed *.sql`)
- Run via dedicated one-shot container at deploy time — **not** at server startup
- Commands:
  - `make migrate-up` — apply all pending migrations
  - `make migrate-down` — roll back most recent migration
  - `make migrate-status` — print current migration state
- Entry point: `backend/cmd/migrate/main.go`

## Query Generation (sqlc)

- Tool: `sqlc v1.30.x` (install: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- Config: `backend/sqlc.yaml`
- Schema source: `backend/migrations/`
- Query source: `backend/db/queries/*.sql`
- Generated output: `backend/db/sqlc/` (committed, do not hand-edit)
- Regenerate: `make sqlc-gen`
- Repository code must only import from `backend/db/sqlc/` — no raw SQL strings in Go

## Local Dev Ports

| Service    | Port |
|------------|------|
| API        | 8080 |
| Frontend   | 3000 |
| PostgreSQL | 5432 |
| Redis      | 6379 |
| MailHog SMTP | 1025 |
| MailHog UI | 8025 |
