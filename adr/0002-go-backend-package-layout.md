# ADR-0002: Go Backend Package Layout and Design-Pattern Conventions

**Status:** Accepted
**Date:** 2026-03-31
**Deciders:** Feather Tech

---

## Context

TASK-002 initialises the Go backend for the e-commerce books platform. Before any feature work begins, the team needs a canonical package layout and a set of design-pattern conventions that all subsequent backend tasks must follow.

The backend is a Go 1.26 / Gin application deployed on AWS. Key constraints:

- **Layered architecture** — the architecture document (§2) mandates a strict Repository → Service → Handler separation. Handlers must never call the database directly.
- **Testability** — each layer must be independently testable via interface injection.
- **Startup safety** — configuration errors (missing env vars) must surface at startup, not at request time.
- **HTTP server configurability** — timeouts, TLS, and future options must be addable without changing the call site in `main.go`.
- **Standard Go project layout** — the `cmd/` / `internal/` split is the community-standard layout for production Go services. Using it ensures tooling (go build, golangci-lint, Docker multi-stage builds) works without custom configuration.

Without explicit conventions, AI coding agents and human contributors will diverge on where code lives, how dependencies flow, and which design patterns to use — making the codebase hard to reason about as it grows.

## Decision

We will use the following package layout and design-pattern conventions for the entire backend:

### Package Layout

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Binary entry point only — wires and starts the server
├── internal/
│   ├── config/              # Adapter pattern: env var → typed Config struct
│   ├── db/                  # Database connection setup (wired in TASK-004)
│   ├── handlers/            # Gin HTTP handlers — parse input, call service, write response
│   ├── middleware/          # Gin middleware (auth, rate-limit, etc.)
│   ├── models/              # Domain structs (no methods that touch the DB)
│   ├── repositories/        # Data-access layer — all SQL/Redis calls live here
│   ├── services/            # Business logic — calls repositories, never gin.Context
│   └── server/              # Builder pattern: constructs *http.Server + Gin router
```

`internal/` prevents external packages from importing the application's private packages, enforcing encapsulation at the module boundary.

### Design Patterns

**Builder — HTTP Server (`internal/server`)**

`ServerBuilder` accumulates configuration via fluent `With*` methods and produces a configured `*http.Server` via `Build()`. This makes `main.go` declarative, decouples the caller from constructor argument order, and allows new options (TLS, custom error handler, additional middleware) to be added without a breaking signature change.

**Adapter — Configuration (`internal/config`)**

`EnvAdapter` is the Adapter that converts raw `os.Getenv` calls into a typed, validated `Config` struct. The adapter calls `log.Fatal` with a descriptive message if any required environment variable is absent, so misconfigured deployments fail immediately at startup rather than silently misbehaving at request time. Future adapters (e.g. AWS Secrets Manager, Vault) can be swapped in by implementing the same interface without changing `main.go`.

**Repository / Service / Handler Layering**

- `handlers` receive `*gin.Context`, validate HTTP input, call one or more `services`, and write HTTP responses. They have no SQL or Redis imports.
- `services` contain business logic. They accept and return domain types from `models`. They have no `*gin.Context` imports, making them independently unit-testable.
- `repositories` own all database interactions. They accept plain Go types and return domain models or errors. Interfaces are defined in the `services` package so that services depend on the abstraction, not the concrete implementation (Dependency Inversion).

### Entry-Point Convention

The binary lives at `cmd/api/main.go`, not `main.go` at the module root. This matches the standard Go multi-binary layout and keeps the Dockerfile `go build ./cmd/api` target stable as new binaries (e.g. `cmd/worker`) are added.

## Consequences

### Positive
- Every future backend task has a clearly defined home for new code — no debate about where a function belongs.
- `internal/` prevents accidental cross-module imports; the compiler enforces the layering.
- The Builder and Adapter patterns isolate change — adding a new config source or server option touches one file, not the call site.
- Services are pure Go with no Gin dependency, making unit tests fast and framework-independent.
- The `cmd/api` entry-point convention supports adding `cmd/worker` (TASK-040) without restructuring.

### Negative
- More directories than a flat layout — adds cognitive overhead for contributors new to standard Go project layout.
- Interface definitions in the `services` package mean touching two files (interface + implementation) when adding a new repository method.

### Neutral
- `doc.go` placeholder files are used in packages that have no implementation yet (e.g. `db/`, `middleware/`) so that `go build ./...` succeeds from day one.
- The `internal/` boundary means the frontend and any future Go CLIs in this monorepo cannot import backend packages directly — any shared types must be placed in a `pkg/` directory if needed in future.

## Alternatives Considered

| Option | Reason rejected |
|--------|----------------|
| Flat layout (`backend/*.go`) | Does not scale past ~5 files; no separation of concerns enforced by the compiler |
| `main.go` at module root | Breaks when a second binary (`cmd/worker`) is needed; non-standard for production Go |
| Functional options for server | Valid alternative to Builder; Builder was chosen for readability — `WithReadTimeout(15s)` is more explicit than passing a slice of `func(*Server)` |
| Config via struct tags + reflection | Adds a third-party dependency; `log.Fatal` on missing vars is simpler and sufficient for this project |
| Interface definitions in `repositories` package | Causes import cycles when services import repositories; placing interfaces in `services` follows the standard "accept interfaces, return structs" Go idiom |

## References

- `docs/02_outputs/03_tech_architecture.md` §2 — Backend API architecture
- `docs/02_outputs/05_implementation_plan.md` TASK-002
- `backend/AGENTS.md` — backend coding conventions
- ADR-0001 — Use Vibecode Framework Structure
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
