---
name: golang
description: Project-specific Go patterns for the e-commerce backend (Gin, layered architecture, sqlc, goose)
paths: "backend/**/*.go"
---

## Project Stack
- Go 1.26, Gin v1.12, PostgreSQL 18, Redis 7.2
- sqlc for generated DB code (`backend/db/sqlc/` — never edit manually)
- goose for migrations (`backend/migrations/`)
- `log/slog` for structured JSON logging (never use `fmt.Println` for logs)

## Package Layout (ADR-0002)
```
backend/
├── cmd/api/          # main.go only — wires everything together
├── internal/
│   ├── config/       # EnvAdapter: reads env → Config struct (fail-fast at startup)
│   ├── db/           # DB connection setup
│   ├── handlers/     # HTTP only: parse input → call service → write response
│   ├── middleware/   # Gin middleware
│   ├── models/       # Pure domain structs (no DB methods, no gin imports)
│   ├── repositories/ # All SQL/Redis — interfaces defined in services package
│   ├── services/     # Business logic — no gin.Context, no DB imports
│   └── server/       # ServerBuilder (builder pattern) → *http.Server
└── db/queries/       # Raw .sql files consumed by sqlc
```

## Hard Rules
- **Handlers never import repositories** — only services
- **Services never import gin** — independently unit-testable
- **Repositories** implement interfaces defined in the `services` package (dependency inversion)
- Config errors → `log.Fatal` at startup, not deferred
- New SQL queries: add to `db/queries/*.sql` then run `make sqlc-gen` — never write raw queries in Go
- New migrations: use goose (`make migrate-new NAME=...`), never alter existing migration files

## Error Handling
```go
// Wrap with context
return fmt.Errorf("createOrder: %w", err)

// Sentinel errors in domain packages
var ErrOrderNotFound = errors.New("order not found")
```

## HTTP Handler Shape
```go
func (h *Handler) CreateOrder(c *gin.Context) {
    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result, err := h.service.CreateOrder(c.Request.Context(), req)
    if err != nil {
        // map domain errors to HTTP status codes here
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
        return
    }
    c.JSON(http.StatusCreated, result)
}
```

## Graceful Shutdown (already wired in main.go)
- SIGINT/SIGTERM triggers 10-second shutdown timeout
- Always propagate `context.Context` to all I/O operations

## Testing
- Services: pure unit tests (no gin, no DB)
- Handlers: use `httptest` with mocked service interfaces
- Repositories: integration tests against real PostgreSQL (never mock the DB)
- Run: `go test ./...` | Race detector: `go test -race ./...`

## Observability
- Use `slog.Info/Error/Warn` with structured key-value pairs
- Handler errors logged at `slog.Error` with `"error"` key
- Request logging handled by Gin middleware (do not duplicate)

## CORS
- Exact origins only — no wildcards (security requirement)
- Configured in `internal/server/` via gin-contrib/cors
