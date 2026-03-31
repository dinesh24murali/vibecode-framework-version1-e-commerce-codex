# Task: TASK-002 — Backend scaffold — Go 1.26 / Gin

**Status:** done
**Created:** 2026-03-31
**ADR refs:** ADR needed — Go backend package layout and design-pattern conventions

---

## Goal

Initialise the Go 1.26 module, install Gin, wire a minimal HTTP server with a `/health` endpoint, and establish the canonical `internal/` package layout that all subsequent backend tasks will follow.

## Background

All backend feature tasks (TASK-005 onward) depend on this scaffold being in place. The package layout must reflect the Repository → Service → Handler layering described in the architecture doc. This task also establishes the design-pattern conventions (Builder for server construction, Adapter for config loading) that must be referenced in the ADR and followed in all future backend tasks.

See `docs/02_outputs/03_tech_architecture.md` §2 (Backend API) and `docs/02_outputs/05_implementation_plan.md` TASK-002.

## Acceptance Criteria

- [ ] `go build ./...` succeeds with zero errors
- [ ] `GET /health` returns `200 OK` with `{"status":"ok"}`
- [ ] Config struct validates all required env vars at startup; server exits with a clear error message if any are missing
- [ ] `golangci-lint run` passes with zero warnings
- [ ] `internal/` package layout is established: `handlers/`, `services/`, `repositories/`, `models/`, `middleware/`, `db/`, `config/`
- [ ] Server construction uses the Builder pattern
- [ ] Config loading from env vars uses the Adapter pattern
- [ ] ADR created in `adr/` documenting the package layout and design-pattern conventions

## Dependencies

- **Tasks:** TASK-001 (monorepo and docker-compose must exist)
- **Memory files:** none (first backend task — memory files will be seeded from TASK-004 onward)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §2, `backend/AGENTS.md`

## Files Expected to Change

- `backend/main.go` *(create)*
- `backend/go.mod` *(create)*
- `backend/go.sum` *(create)*
- `backend/internal/config/config.go` *(create)*
- `backend/internal/server/server.go` *(create)*
- `backend/internal/server/router.go` *(create)*
- `backend/Makefile` *(create)*
- `adr/NNNN-go-backend-package-layout.md` *(create)*

## Notes

- Go version: `1.26`; Gin version: `1.12.0` (per architecture doc)
- The server entry point is `backend/main.go` (or `backend/cmd/api/main.go` — decide in scratch file and note in ADR)
- Do not wire Postgres or Redis connections here — that is TASK-004
- The `/health` endpoint must be reachable without auth middleware
- `golangci-lint` config should be added at `backend/.golangci.yaml` if not already present at repo root

---

## Scratch File

AI: before implementing, create `tasks/active/task-002-backend-scaffold-go-gin.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
