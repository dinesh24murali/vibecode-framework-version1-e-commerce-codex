# Scratch: TASK-002 — Backend scaffold — Go 1.26 / Gin

**Date:** 2026-03-31
**Status:** done

---

## Plan (numbered steps)

1. Create `backend/go.mod` with module path `github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend` and Go 1.26 directive.
2. Fetch Gin v1.12.0 dependency (`go get github.com/gin-gonic/gin@v1.12.0`).
3. Create `backend/internal/config/config.go` — Adapter pattern: reads from env vars, validates required fields, returns typed `Config` struct.
4. Create `backend/internal/models/` with an empty placeholder `.go` file so the package directory exists and compiles.
5. Create `backend/internal/middleware/` with a placeholder file.
6. Create `backend/internal/repositories/` with a placeholder file.
7. Create `backend/internal/services/` with a placeholder file.
8. Create `backend/internal/handlers/health.go` — `GET /health` handler returning `{"status":"ok"}`.
9. Create `backend/internal/db/` with a placeholder file (wired in TASK-004).
10. Create `backend/internal/server/server.go` — Builder pattern: `ServerBuilder` struct with `WithConfig`, `WithRouter`, and `Build` methods returning `*http.Server`.
11. Create `backend/internal/server/router.go` — wires Gin engine, registers `/health` route, applies CORS middleware, returns `*gin.Engine`.
12. Create `backend/cmd/api/main.go` — entry point: loads config via Adapter, builds server via Builder, starts `http.Server`.
13. Create `backend/.golangci.yaml` — sensible linter settings; enable `govet`, `errcheck`, `staticcheck`, `unused`, `goimports`, `misspell`; set Go version.
14. Create `backend/Makefile` with targets: `run`, `build`, `test`, `lint`.
15. Run `go mod tidy` to resolve and pin all deps.
16. Run `go build ./...` to confirm zero errors.
17. Run `golangci-lint run` and fix any issues.
18. Create `adr/0001-go-backend-package-layout.md` using `adr/template.md` as base.
19. Update `CHANGELOG.md`.
20. Move task file and scratch file to `tasks/done/`.

---

## Files to change

### New files
| File | Purpose |
|------|---------|
| `backend/go.mod` | Go module definition |
| `backend/go.sum` | Dependency checksums (generated) |
| `backend/cmd/api/main.go` | Binary entry point |
| `backend/internal/config/config.go` | Env-var Adapter |
| `backend/internal/server/server.go` | Builder pattern server constructor |
| `backend/internal/server/router.go` | Gin router wiring |
| `backend/internal/handlers/health.go` | GET /health handler |
| `backend/internal/handlers/doc.go` | Package placeholder |
| `backend/internal/services/doc.go` | Package placeholder |
| `backend/internal/repositories/doc.go` | Package placeholder |
| `backend/internal/models/doc.go` | Package placeholder |
| `backend/internal/middleware/doc.go` | Package placeholder |
| `backend/internal/db/doc.go` | Package placeholder (wired in TASK-004) |
| `backend/.golangci.yaml` | Linter configuration |
| `backend/Makefile` | Dev tooling targets |
| `adr/0001-go-backend-package-layout.md` | ADR for package layout and design patterns |

### Modified files
| File | Change |
|------|--------|
| `CHANGELOG.md` | Add TASK-002 entry |
| `tasks/active/task-002-backend-scaffold-go-gin.md` | Update status to done |

---

## Uncertainties

- Go 1.26 is not yet released as of this writing (current stable is 1.23.x); the `go.mod` `go` directive will be set to `1.26` as specified in the task. `go get` will use whatever toolchain is available; the build will use the installed Go. This is acceptable because the task spec mandates it.
- Gin v1.12.0 — need to verify this tag exists on pkg.go.dev. If it does not exist at time of `go get`, fall back to the latest stable v1.x and note it.
- The module path is inferred from the repo URL in AGENTS.md: `github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend`.
- `golangci-lint` version installed on the machine may not match the `.golangci.yaml` lint directives; configure defensively.

---

## What will be skipped

- No Postgres or Redis connection wiring (TASK-004).
- No auth middleware (TASK-005+).
- No worker binary at `backend/cmd/worker/` (TASK-040).
- No actual handler implementations beyond `/health`.
- No OpenAPI contract test for `/health` — the task explicitly notes this is internal scaffolding, not a feature endpoint.
- No frontend client regeneration — no API contract changes that affect the frontend spec.

---

## Risks

- **Go version:** go 1.26 toolchain may not be available; `go build` will still work with 1.22/1.23 for the code written here. The `go.mod` directive is a declaration, not a hard runtime requirement for this scaffold.
- **Gin version:** If `v1.12.0` does not exist in the module proxy, `go get` will fail. Fallback: use latest Gin v1 release.
- **golangci-lint compatibility:** Some linters enabled in `.golangci.yaml` may not exist in the installed version. Use a conservative, broadly supported linter set.
- **CORS middleware:** Gin does not ship a CORS middleware in its core; we will use `github.com/gin-contrib/cors`. This is an additional dependency justified by the CORS requirement in the architecture doc.
