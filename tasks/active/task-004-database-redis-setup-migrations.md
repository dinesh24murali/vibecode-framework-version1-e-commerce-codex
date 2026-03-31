# Task: TASK-004 — Database + Redis setup, migrations & query generation tooling

**Status:** active
**Created:** 2026-03-31
**ADR refs:** ADR needed — SQL-first persistence with sqlc and goose (vs ORM)

---

## Goal

Set up three foundational infrastructure layers in the backend: Postgres connection pool (pgxpool), goose migration tooling, Redis client, and sqlc query generation — all wired so subsequent tasks can build on them immediately.

## Background

Every backend feature task from TASK-005 onward depends on these three layers being operational. This task does not write the business schema (that is TASK-005) — it installs and configures the tooling and connections. Migrations must not auto-run on server startup; they are run explicitly via `make migrate-up`.

See `docs/02_outputs/03_tech_architecture.md` §2 (Database), §2 (Cache and Ephemeral State), §3 (Migration Strategy) and `docs/02_outputs/05_implementation_plan.md` TASK-004.

## Acceptance Criteria

- [ ] `make migrate-up` runs all migrations against the local Postgres instance without error
- [ ] `make migrate-down` rolls back the most recent migration
- [ ] `make migrate-status` prints the current migration state
- [ ] `make sqlc-gen` runs `sqlc generate` without errors (requires schema to exist first)
- [ ] `sqlc.yaml` committed at repo root; `backend/db/sqlc/` generated files committed (or excluded with a documented `.gitignore` note)
- [ ] `pgxpool` connection honours configured limits: `MaxConns: 20`, `MinConns: 4`, `MaxConnLifetime: 30m`, `MaxConnIdleTime: 5m`
- [ ] Redis client connects and pings successfully on startup
- [ ] `go test ./internal/db/...` and `go test ./internal/cache/...` verify connection state
- [ ] Server logs a clear error and exits non-zero if `DATABASE_URL` or `REDIS_URL` is unreachable

## Dependencies

- **Tasks:** TASK-001 (docker-compose with postgres and redis must exist), TASK-002 (Go module and server scaffold must exist)
- **Memory files:** none (will create `memory/api.md` and `memory/infra.md` entries after this task)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §2 (Database), §2 (Cache), §3 (Migration Strategy), §9 (DB Connection Pooling)

## Files Expected to Change

- `backend/internal/db/db.go` *(create — pgxpool setup)*
- `backend/internal/db/migrate.go` *(create — goose integration with embedded migrations)*
- `backend/internal/cache/redis.go` *(create — go-redis/v9 client)*
- `backend/migrations/` *(create directory — goose migration files go here)*
- `backend/db/queries/` *(create directory — raw .sql query files for sqlc)*
- `backend/db/sqlc/` *(create directory — sqlc-generated Go code, do not hand-edit)*
- `sqlc.yaml` *(create at repo root or `backend/` — points schema at `migrations/`, queries at `db/queries/`, output to `db/sqlc/`)*
- `backend/Makefile` *(update with migrate-up, migrate-down, migrate-status, sqlc-gen targets)*
- `memory/infra.md` *(create or update after task — document pool limits, Redis AOF requirement)*

## Notes

- Library versions: `pgx v5.9.x`, `pressly/goose v3.27.x`, `go-redis/v9`, `sqlc v1.30.x`
- `migrate.go` must embed the `migrations/` directory using Go's `embed` package so migrations are bundled into the binary — migrations are run via a dedicated one-shot container at deploy time, not at server startup
- sqlc `schema` path points to `backend/migrations/`; queries path points to `backend/db/queries/`; output package is `backend/db/sqlc/`
- Repository implementations in `internal/repositories/` must import only from `backend/db/sqlc/` — no raw SQL strings in Go code
- Redis AOF persistence must be enabled in `docker-compose.yml` because Redis also carries the job queue (per architecture doc §2 Cache)
- This task creates the ADR for SQL-first persistence (sqlc + goose vs ORM)

---

## Scratch File

AI: before implementing, create `tasks/active/task-004-database-redis-setup-migrations.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
