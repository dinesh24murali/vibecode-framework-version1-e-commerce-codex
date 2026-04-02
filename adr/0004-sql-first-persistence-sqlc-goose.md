# ADR-0004: SQL-first persistence with sqlc and goose (vs ORM)

**Status:** Accepted
**Date:** 2026-04-01
**Deciders:** team

---

## Context

The backend needs a persistence strategy for a PostgreSQL database. The two main
schools of thought are:

1. **ORM** (e.g. GORM, ent) — Go structs drive schema and queries; SQL is generated automatically.
2. **SQL-first** — raw SQL migrations define schema; type-safe Go code is generated from hand-written queries.

This project prioritises correctness, auditability, and predictable query plans over
developer convenience. The business domain (inventory, orders, payments) demands fine-grained
control over queries, indexes, and transactions.

## Decision

We will use **sqlc** for query code generation and **goose** for migrations.

- Schema is defined in `backend/migrations/*.sql` files managed by goose.
- Queries are written by hand in `backend/db/queries/*.sql`.
- sqlc generates type-safe Go code in `backend/db/sqlc/` from those queries.
- Repository implementations in `internal/repositories/` import only from `db/sqlc/` — no raw SQL strings in Go code.
- Migrations run via `make migrate-up` (one-shot container at deploy time); they do **not** auto-run on server startup.
- Migration files are embedded in the binary via `backend/migrations/embed.go` using Go's `embed` package.

## Consequences

### Positive
- Every executed query is visible in version control as plain SQL — easy to review and audit.
- sqlc eliminates hand-written `Scan` boilerplate and catches query/struct mismatches at code-gen time.
- goose provides explicit migration versioning and rollback support.
- No ORM reflection overhead at runtime.

### Negative
- Developers must write SQL manually; no auto-generated CRUD.
- Schema changes require both a migration file and updated query files; more files to coordinate than an ORM.
- sqlc codegen must be re-run after every query or schema change (`make sqlc-gen`).

### Neutral
- `database/sql` interface is used by goose (via pgx stdlib driver); pgxpool is used for application queries — two slightly different interfaces in the same project.

## Alternatives Considered

| Option | Reason rejected |
|--------|----------------|
| GORM | Hides SQL; generated queries are unpredictable; poor support for complex joins and CTEs |
| ent (Facebook) | Schema-as-code is powerful but adds significant framework lock-in and a steep learning curve |
| Plain `database/sql` with hand-written scans | No codegen safety net; high boilerplate; error-prone |
| bun ORM | Lighter than GORM but still ORM semantics; sqlc gives stronger compile-time guarantees |

## References

- `tasks/active/task-004-database-redis-setup-migrations.md`
- `docs/02_outputs/03_tech_architecture.md` §2 (Database), §3 (Migration Strategy)
- [sqlc documentation](https://docs.sqlc.dev)
- [goose documentation](https://pressly.github.io/goose)
