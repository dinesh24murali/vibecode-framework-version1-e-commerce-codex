# e-commerce-site

A books-focused e-commerce platform for selling books online.

**Stack:** Go 1.26 / Gin · NextJS 16 SSG · PostgreSQL 18 · Redis 7.2 · JWT (Ed25519) · shadcn/ui · AWS

---

## Prerequisites

- Docker Engine 29.x + Docker Compose v2
- Go 1.26+
- Node.js 22.x + npm
- `sqlc` v1.30.x (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- `golangci-lint` (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

---

## Local Development Setup

```bash
# 1. Clone the repository
git clone https://github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex
cd vibecode-framework-version1-e-commerce-codex

# 2. Copy environment file and fill in values
cp .env.example .env
# Edit .env — at minimum set JWT_PRIVATE_KEY, JWT_PUBLIC_KEY

# 3. Start all local services (Postgres, Redis, MailHog, api, worker, frontend)
docker compose up

# 4. Run database migrations
make migrate-up

# 5. (Optional) Seed development data
make seed
```

### Running without Docker

```bash
# Terminal 1: start backend + frontend dev servers concurrently
make dev

# Or individually:
cd backend && go run ./cmd/api      # API on :8080
cd frontend && npm run dev          # Frontend on :3000
```

---

## Makefile Targets

| Command | Description |
|---------|-------------|
| `make dev` | Start backend and frontend dev servers concurrently |
| `make build` | Build all Docker images |
| `make test` | Run backend unit and integration tests |
| `make verify` | Run e2e + DOM verification checks |
| `make migrate-up` | Apply all pending database migrations |
| `make migrate-down` | Roll back the most recent migration |
| `make migrate-status` | Show current migration state |
| `make seed` | Seed the database with development fixture data |
| `make sqlc-gen` | Regenerate Go database code from SQL files |
| `make codegen` | Regenerate frontend TypeScript API client from OpenAPI spec |
| `make adr SLUG=my-decision` | Create a new Architecture Decision Record |
| `make task NAME=my-feature` | Create a new task file from the template |
| `make check-env` | Validate `.env` has all required variables |
| `make help` | Show all available commands |

---

## Local Service URLs

| Service | URL |
|---------|-----|
| Frontend (storefront) | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| API health check | http://localhost:8080/health |
| MailHog (email UI) | http://localhost:8025 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

---

## Project Structure

```
.
├── backend/                # Go 1.26 API + worker
│   ├── cmd/api/            # API binary entry point
│   ├── cmd/worker/         # Worker binary entry point
│   ├── cmd/migrate/        # Migration runner entry point
│   ├── internal/           # Application code (handlers, services, repositories)
│   ├── migrations/         # goose SQL migration files
│   └── db/
│       ├── queries/        # Raw .sql query files (consumed by sqlc)
│       └── sqlc/           # sqlc-generated Go code (do not hand-edit)
├── frontend/               # NextJS 16 SSG storefront
│   ├── app/                # App Router pages and layouts
│   ├── components/         # Reusable UI components
│   ├── lib/
│   │   ├── api/            # Generated API client (from OpenAPI spec)
│   │   └── store/          # Zustand v5 state slices
│   └── public/
├── docs/
│   ├── 00_intake/          # Project intake questionnaire
│   ├── 01_prompts/         # AI prompt templates
│   └── 02_outputs/         # AI-generated documents (PRD, arch, API spec, etc.)
├── adr/                    # Architecture Decision Records
├── memory/                 # AI agent memory files (domain knowledge)
├── tasks/
│   ├── active/             # In-progress tasks + scratch files
│   └── done/               # Completed tasks
├── tests/
│   ├── e2e/                # Playwright end-to-end tests
│   └── contract/           # OpenAPI contract tests
├── verify/                 # Verification scripts and playbooks
└── AGENTS.md               # Master AI instructions
```

---

## AI Agent Instructions

All AI tool instructions are in `AGENTS.md`. `CLAUDE.md` and `.cursorrules` both redirect there.

For AI-assisted development, follow the two-phase workflow in `AGENTS.md`:
1. **Before coding:** create `tasks/active/<task-name>.scratch.md` with a plan
2. **After coding:** update `CHANGELOG.md`, move task to `tasks/done/`, update `memory/`

---

## Production

Production runs on a single AWS EC2 instance with Docker Compose. See `docs/02_outputs/03_tech_architecture.md` for the full architecture.

- Frontend: S3 + CloudFront (static export)
- Backend: EC2 running Docker Compose (api, worker, postgres, redis)
- Payments: Razorpay
- Email: Amazon SES
- Secrets: SOPS-encrypted `.env` files + AWS SSM Parameter Store
