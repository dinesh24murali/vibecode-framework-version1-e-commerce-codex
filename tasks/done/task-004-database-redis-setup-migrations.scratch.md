# Scratch: TASK-004 — Database + Redis setup, migrations & query generation tooling

## Plan

1. **Add Go dependencies** — `pgx/v5`, `pressly/goose/v3`, `go-redis/v9` via `go get` in `backend/`
2. **Create `backend/migrations/embed.go`** — exports `var FS embed.FS` with `//go:embed *.sql`; this is the embed anchor co-located with migration files
3. **Create `backend/internal/db/db.go`** — `Open(ctx, url)` returning `*pgxpool.Pool`; honours MaxConns:20, MinConns:4, MaxConnLifetime:30m, MaxConnIdleTime:5m; pings on startup
4. **Create `backend/internal/db/migrate.go`** — imports `migrations` package for the embedded FS; exports `RunUp`, `RunDown`, `PrintStatus` wrapping goose
5. **Create `backend/internal/cache/redis.go`** — exports `NewClient(url)` returning `(*redis.Client, error)`; pings on startup
6. **Create `backend/migrations/20260401000000_init.sql`** — empty UP/DOWN so goose has a valid starting point
7. **Create `backend/db/queries/` and `backend/db/sqlc/`** — placeholder dirs
8. **Create `backend/sqlc.yaml`** — schema: `migrations/`, queries: `db/queries/`, output: `db/sqlc/`
9. **Update `backend/Makefile`** — add `migrate-up`, `migrate-down`, `migrate-status`, `sqlc-gen` targets
10. **Wire `main.go`** — open pgxpool and Redis after config load; `log.Fatal` + exit if either fails
11. **Update `backend/internal/server/server.go`** — accept `*pgxpool.Pool` and `*redis.Client` as fields for future handler DI
12. **Write integration tests** — `backend/internal/db/db_test.go` and `backend/internal/cache/redis_test.go`; skip if env vars absent
13. **Create `adr/0004-sql-first-persistence-sqlc-goose.md`**
14. **Create `memory/infra.md`**

## Files to Change / Create

| File | Action |
|------|--------|
| `backend/go.mod` / `go.sum` | updated by `go get` |
| `backend/migrations/embed.go` | create — embed anchor |
| `backend/migrations/20260401000000_init.sql` | create — empty seed migration |
| `backend/internal/db/db.go` | create |
| `backend/internal/db/db_test.go` | create |
| `backend/internal/db/migrate.go` | create |
| `backend/internal/db/doc.go` | update comment |
| `backend/internal/cache/redis.go` | create |
| `backend/internal/cache/redis_test.go` | create |
| `backend/db/queries/.gitkeep` | create |
| `backend/db/sqlc/.gitkeep` | create |
| `backend/sqlc.yaml` | create |
| `backend/Makefile` | update |
| `backend/cmd/api/main.go` | update — wire DB + Redis |
| `backend/internal/server/server.go` | update — accept pool + redis |
| `adr/0004-sql-first-persistence-sqlc-goose.md` | create |
| `memory/infra.md` | create |

## Uncertainties

- **`//go:embed` path constraint**: The embed directive file must be in a directory that is an ancestor of (or the same as) the embedded files. Since migrations live at `backend/migrations/`, the embed must be declared in that same directory — hence `backend/migrations/embed.go`. `migrate.go` in `internal/db/` will import the `migrations` package to get the `FS`.
- **Server builder DI**: Currently `ServerBuilder` has no DB/Redis fields. Adding them now is the right call so TASK-005 can wire handlers cleanly, but it means updating the builder and `Build()` signature. Keeping changes minimal — add fields + `WithPool` / `WithRedis` fluent methods only.
- **sqlc with empty queries dir**: `sqlc generate` will succeed but produce nothing — acceptable for this task.

## What Will Be Skipped

- Business schema SQL — TASK-005
- Real `.sql` query files — come with TASK-005+
- Full handler DI of pool — TASK-005/006

## Risks

- `go get` may pull indirect deps that interact with existing ones; will verify `go build ./...` passes after
- Integration tests require running Postgres + Redis; they skip gracefully if env vars are absent so CI won't break
