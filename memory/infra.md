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

## Docker-only toolchain (task 000)
- All go/node/sqlc/goose/lint/playwright/terraform commands run via `docker-compose.tools.yml` (services go, node, e2e, static, tf; profile `tools`). Inside containers use service hostnames, not localhost.
- Postgres 18 volume is `postgres18_data` mounted at `/var/lib/postgresql` (the old `postgres_data` volume is unused).
- `make` is not installed on the Windows host; run the raw `docker compose -f docker-compose.yml -f docker-compose.tools.yml run --rm <svc> <cmd>` form with `MSYS_NO_PATHCONV=1`.
- Frontend static export is served by the `static` busybox service at http://static:3000; Playwright deps install to /opt/tools in the e2e container.

## Background Worker (TASK-040)

- Binary `backend/cmd/worker`; reads `DATABASE_URL`/`REDIS_URL` only (no JWT keys needed); exits non-zero if either is unreachable.
- Stream `jobs`, consumer group `workers`, dead-letter stream `jobs:dead-letter`. Job fields: `id`, `type`, `payload`, `attempts`.
- Failure: retry up to 3 times (backoff 1s doubling, cap 60s), then dead-letter. Stale pending entries reclaimed with XAUTOCLAIM (1 min idle).
- Idempotency: Redis key `job:done:<id>` (7d TTL) skips replays. Register handlers via `worker.Dispatcher.Register`; Phase 2 stubs live in `internal/worker/handlers/stubs.go`.

## Image uploads (TASK-013)

- `POST /api/v1/admin/upload/image/presign` (admin only): returns a presigned S3 PUT (5 min, content type and length signed) plus the CDN `public_url` that gets stored on products. Env: `S3_BUCKET`, `AWS_REGION`, `CDN_BASE_URL`, standard AWS credentials. Empty `S3_BUCKET` disables it (503). `auth.RequireRole` guards admin routes (403 for non-admin).
