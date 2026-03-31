# Developer Setup Guide — e-commerce-site

> Target: a brand-new engineer should have a fully working local environment in under 30 minutes.

---

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Repository Setup](#2-repository-setup)
3. [Local Development](#3-local-development)
4. [Database Setup](#4-database-setup)
5. [Running Tests](#5-running-tests)
6. [Common Development Tasks](#6-common-development-tasks)
7. [Debugging](#7-debugging)
8. [CI/CD Overview](#8-cicd-overview)
9. [Deployment (Non-Prod)](#9-deployment-non-prod)

---

## 1. Prerequisites

Install the following tools before cloning the repository. Versions listed are minimums.

| Tool | Min Version | Install |
|------|-------------|---------|
| Git | 2.40 | https://git-scm.com/downloads |
| Docker Desktop (macOS/Windows) | 4.30 | https://www.docker.com/products/docker-desktop |
| Docker Engine + Compose v2 (Linux) | 29.x / 2.27 | https://docs.docker.com/engine/install |
| Go | 1.26 | https://go.dev/dl |
| Node.js | 22.x | https://nodejs.org/en/download (use `nvm` recommended) |
| Make | 3.81 | Pre-installed on macOS/Linux. Windows: https://gnuwin32.sourceforge.net/packages/make.htm |
| goose (migrations) | 3.27 | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| sqlc (query codegen) | 1.30 | https://docs.sqlc.dev/en/latest/overview/install.html |

**macOS** — use Homebrew for convenience:

```bash
brew install go node make
brew install --cask docker
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**Linux (Ubuntu/Debian)**:

```bash
sudo apt-get update && sudo apt-get install -y make git curl
# Install Docker Engine
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER   # then log out and back in
# Install Go 1.26 from https://go.dev/dl/
# Install Node.js 22 via nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
nvm install 22 && nvm use 22
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**Windows** — use WSL2 with Ubuntu and follow the Linux instructions above. Native Windows is not supported.

> **Warning:** Docker Desktop on macOS with Apple Silicon (M1/M2/M3) is fully supported. No extra steps are needed; all images build for `linux/arm64` automatically via Compose.

---

## 2. Repository Setup

### Clone the repository

```bash
git clone https://github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex.git
cd vibecode-framework-version1-e-commerce-codex
```

### Copy and fill environment variables

The project uses a `.env.example` file at the repository root as the canonical list of all required variables.

```bash
cp .env.example .env
```

Open `.env` and fill in every value. The table below describes each variable and where to obtain it.

| Variable | Description | How to obtain |
|----------|-------------|---------------|
| `POSTGRES_USER` | DB superuser name | Pick any value, e.g. `ecomm` |
| `POSTGRES_PASSWORD` | DB superuser password | Pick any strong value |
| `POSTGRES_DB` | Database name | Pick any value, e.g. `ecomm_dev` |
| `POSTGRES_HOST` | DB hostname (internal) | `postgres` (Docker service name) |
| `REDIS_HOST` | Redis hostname (internal) | `redis` |
| `JWT_ED25519_PRIVATE_KEY` | Ed25519 private key (base64) | Generate locally — see below |
| `JWT_ACCESS_TTL_SECONDS` | Access token TTL | `900` (15 min) for dev |
| `JWT_REFRESH_TTL_SECONDS` | Refresh token TTL | `604800` (7 days) for dev |
| `RAZORPAY_KEY_ID` | Razorpay API key | Razorpay dashboard → API Keys |
| `RAZORPAY_KEY_SECRET` | Razorpay API secret | Razorpay dashboard → API Keys |
| `RAZORPAY_WEBHOOK_SECRET` | Webhook signing secret | Razorpay dashboard → Webhooks |
| `MAILHOG_HOST` | Email relay host | `mailhog` (Docker service) |
| `MAILHOG_PORT` | Email relay port | `1025` |
| `APP_ENV` | Environment name | `dev` |
| `LOG_LEVEL` | Log verbosity | `debug` for local dev |

### Generate JWT signing keys locally

```bash
# Generate Ed25519 key pair (requires OpenSSL 3.x)
openssl genpkey -algorithm ed25519 -out ed25519.pem
openssl pkey -in ed25519.pem -pubout -out ed25519.pub.pem

# Base64-encode for .env
echo "JWT_ED25519_PRIVATE_KEY=$(base64 -w 0 ed25519.pem)"
echo "JWT_ED25519_PUBLIC_KEY=$(base64 -w 0 ed25519.pub.pem)"

# Clean up raw key files
rm ed25519.pem ed25519.pub.pem
```

Copy the output values into `.env`.

> **Warning:** Never commit `.env` to git. It is already in `.gitignore`. Never share JWT keys or Razorpay secrets in Slack, email, or PR comments.

### Razorpay credentials (dev)

For local development, use Razorpay **test mode** keys. Log in to the [Razorpay Dashboard](https://dashboard.razorpay.com), switch to **Test Mode**, and navigate to **Settings → API Keys**. Generate a key pair and copy both values into `.env`.

---

## 3. Local Development

All services run inside Docker Compose. You do not need Postgres or Redis installed on your host machine.

### Start all services

```bash
docker compose up --build
```

On first run, Docker will pull base images and build the `api`, `worker`, and `frontend` images. This takes 3–5 minutes. Subsequent starts are fast.

### Expected healthy output

```
postgres   | database system is ready to accept connections
redis      | Ready to accept connections tcp
api        | {"level":"info","service":"api","msg":"server started","addr":":8080"}
worker     | {"level":"info","service":"worker","msg":"worker started"}
frontend   | ▲ Next.js 16.x
frontend   | Local: http://localhost:3000
mailhog    | [APIv1] KEEPALIVE /api/v2/events
```

| Service | URL |
|---------|-----|
| Storefront | http://localhost:3000 |
| Backend API | http://localhost:8080/api/v1 |
| API health check | http://localhost:8080/healthz |
| MailHog web UI | http://localhost:8025 |
| PostgreSQL | `localhost:5432` (connect with any Postgres client) |
| Redis | `localhost:6379` |

### Stop all services

```bash
docker compose down
```

To also remove volumes (wipes the database):

```bash
docker compose down -v
```

### Run only the backend (no frontend)

```bash
docker compose up postgres redis mailhog api worker
```

Or run Go directly on your host (fastest for backend iteration):

```bash
cd backend
go run ./cmd/server
```

Ensure `.env` is exported in your shell, or use a tool like `direnv`.

### Run only the frontend (no backend)

```bash
cd frontend
npm install
npm run dev
```

The frontend dev server starts at http://localhost:3000 and proxies API calls to `http://localhost:8080` (configured in `next.config.js`).

---

## 4. Database Setup

### Run migrations

Migrations use `goose` and live in `backend/migrations/`.

```bash
# From the repo root — applies all pending migrations
make migrate-up

# Equivalent manual command:
goose -dir backend/migrations postgres \
  "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" up
```

> **Warning:** Migrations run on the *host-exposed* port `5432`. Make sure the `postgres` Docker service is running before running migrations from your host.

### Check migration status

```bash
make migrate-status
# or
goose -dir backend/migrations postgres \
  "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" status
```

### Seed development data

```bash
make seed
# or
go run ./backend/cmd/seed
```

This inserts a default admin user, sample book categories, and a small product catalog.

Default admin credentials (dev only):

| Field | Value |
|-------|-------|
| Email | `admin@example.com` |
| Password | `devpassword123` |

### Reset the database

```bash
# Destroys all data and re-runs migrations + seed
make db-reset

# Equivalent steps:
docker compose down -v postgres
docker compose up -d postgres
make migrate-up
make seed
```

---

## 5. Running Tests

### Unit tests (backend)

```bash
cd backend
go test ./...
```

With coverage output:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html   # macOS
xdg-open coverage.html  # Linux
```

### Unit tests (frontend)

```bash
cd frontend
npm test
```

### Integration tests (backend)

Integration tests require a running Postgres instance. The test suite uses a real test database — **do not mock the database**.

```bash
# Start only the database dependency
docker compose up -d postgres redis

cd backend
go test -tags integration ./...
```

The test runner creates a fresh schema in the `ecomm_test` database automatically.

### E2E tests

```bash
make verify
```

This runs two steps in sequence:

```bash
bash verify/scripts/run-e2e.sh        # Playwright e2e tests
npx ts-node verify/scripts/check-dom.ts  # DOM structure assertions
```

All services must be running (`docker compose up`) before running e2e tests.

To run a single Playwright test file:

```bash
cd tests/e2e
npx playwright test checkout.spec.ts
```

To open the Playwright UI for interactive debugging:

```bash
cd tests/e2e
npx playwright test --ui
```

### Coverage gate

Coverage must not fall below the thresholds in `tests/coverage-baseline.json`. To verify:

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

Compare the total against the threshold in `tests/coverage-baseline.json`. If you add a new package, add its threshold to that file.

---

## 6. Common Development Tasks

### Add a new API endpoint

> **Warning:** Read `memory/api.md` and `memory/auth.md` before modifying any handler. All new endpoints must be documented in `docs/02_outputs/04_api_spec.yaml` before implementation.

1. Create a task file and scratch file first (see below).
2. Add the endpoint to `docs/02_outputs/04_api_spec.yaml`.
3. Regenerate `sqlc` types if new queries are needed:
   ```bash
   cd backend && sqlc generate
   ```
4. Add the handler in `backend/internal/handlers/`.
5. Add business logic in `backend/internal/services/`.
6. Register the route in `backend/internal/router/`.
7. Add a contract test in `tests/contract/`.
8. Run backend tests:
   ```bash
   cd backend && go test ./...
   ```

### Add a new frontend page/route

> **Warning:** Read `memory/auth.md` before touching auth-related components. Never write API client code manually — regenerate from the OpenAPI spec.

1. Create a task file and scratch file first.
2. If the API spec changed, regenerate the frontend client:
   ```bash
   # Run from repo root — replace with the actual codegen command once configured
   make codegen-frontend
   ```
3. Add the page in `frontend/src/app/<route>/page.tsx`.
4. Add an e2e test in `tests/e2e/<feature>.spec.ts`.
5. Run verification:
   ```bash
   make verify
   ```

### Create a new database migration

```bash
# Creates a new timestamped migration file
make migration NAME=add_product_tags
# or
goose -dir backend/migrations create add_product_tags sql
```

Edit the generated file in `backend/migrations/`. Then apply:

```bash
make migrate-up
```

Migration rules:
- Forward-only in production
- Use expand-and-contract for destructive changes
- Keep seed/data migrations separate from structural schema changes

### Create a new ADR

```bash
make adr SLUG=use-redis-streams-for-queue
```

This copies `adr/template.md` to `adr/NNNN-use-redis-streams-for-queue.md`. Fill in the context, decision, and consequences. Set status to `Accepted` once the decision is final.

Do not modify an `Accepted` ADR. If a decision changes, create a new ADR that supersedes the old one and update its status field.

### Create a new task

```bash
make task NAME=add-coupon-endpoint
```

This creates `tasks/active/add-coupon-endpoint.md` from `tasks/_template.md`.

**Before writing any code**, create a scratch file:

```bash
# Create manually or with:
touch tasks/active/add-coupon-endpoint.scratch.md
```

The scratch file must contain:
1. Step-by-step implementation plan
2. Every file that will be created or modified
3. Open questions and uncertainties
4. Out-of-scope items
5. Risks

Once the task is done:

```bash
mv tasks/active/add-coupon-endpoint.md tasks/done/
mv tasks/active/add-coupon-endpoint.scratch.md tasks/done/
# Update CHANGELOG.md and relevant memory/ files
```

---

## 7. Debugging

### Enable debug logging

In `.env`, set:

```bash
LOG_LEVEL=debug
APP_ENV=dev
```

Restart the affected service:

```bash
docker compose restart api
```

All application logs are structured JSON. To pretty-print in the terminal:

```bash
docker compose logs -f api | jq .
```

### Connect a debugger to the backend

Run the backend outside Docker for debugger attachment:

```bash
cd backend
# Build with debug symbols (no optimization)
go build -gcflags="all=-N -l" -o bin/server ./cmd/server
./bin/server
```

**VS Code** — add to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to Go server",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickProcess}"
    }
  ]
}
```

**Delve (CLI)**:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/server
```

### Connect a debugger to the frontend

The Next.js dev server (`npm run dev`) exposes Node.js inspector on port `9229`.

**VS Code** — add to `.vscode/launch.json`:

```json
{
  "name": "Next.js: debug server-side",
  "type": "node",
  "request": "attach",
  "port": 9229,
  "skipFiles": ["<node_internals>/**"]
}
```

For client-side debugging, use the browser DevTools directly.

### Common errors and solutions

| Error | Cause | Fix |
|-------|-------|-----|
| `dial tcp 127.0.0.1:5432: connect: connection refused` | Postgres not running | `docker compose up -d postgres` |
| `pq: password authentication failed` | Wrong DB password in `.env` | Check `POSTGRES_PASSWORD` matches Compose service |
| `jwt: token is expired` | Access token TTL too short in dev | Increase `JWT_ACCESS_TTL_SECONDS` in `.env` |
| `goose: no migration files found` | Wrong `-dir` path | Run from repo root with `make migrate-up` |
| `EADDRINUSE :::3000` | Port 3000 already in use | `lsof -i :3000` then kill the process |
| `docker: Cannot connect to the Docker daemon` | Docker not running | Start Docker Desktop / `sudo systemctl start docker` |
| `permission denied while trying to connect to Docker` | User not in docker group (Linux) | `sudo usermod -aG docker $USER` then log out/in |
| Frontend shows stale data | API client out of sync with spec | Regenerate: `make codegen-frontend` |
| `sqlc generate` fails | Query/schema mismatch | Ensure latest migration is applied before running codegen |

---

## 8. CI/CD Overview

CI/CD is not yet configured for this project. All pipeline steps are currently performed manually.

The intended future pipeline (GitHub Actions) will:

**On pull request:**
- Run `go test ./...` (backend unit + integration tests)
- Run `npm test` (frontend unit tests)
- Run `golangci-lint run` (backend lint)
- Build Docker images (without pushing)
- Run `make verify` (e2e tests against a disposable Compose environment)
- Report coverage and fail if thresholds in `tests/coverage-baseline.json` are not met

**On merge to `main`:**
- All PR checks above, plus:
- Build and tag versioned Docker images for `api` and `worker`
- Build the Next.js static export
- Upload the static export to the S3 bucket
- Invalidate the relevant CloudFront paths
- Push images to the container registry
- Trigger deployment to the EC2 host

**To check pipeline status** (once configured):

```bash
gh run list --repo dinesh24murali/vibecode-framework-version1-e-commerce-codex
gh run view <run-id>
```

---

## 9. Deployment (Non-Prod)

There is no automated staging environment. When a manual deployment is needed (e.g. for stakeholder review), follow these steps on the target EC2 host.

> **Warning:** These steps mutate the production-equivalent environment. Confirm with the team before proceeding.

### One-time host setup

```bash
# On the EC2 host
sudo apt-get install -y docker.io docker-compose-plugin
sudo usermod -aG docker ubuntu
```

### Manual deployment steps

```bash
# 1. Build images locally
docker build -t ecomm-api:latest ./backend
docker build -t ecomm-worker:latest ./backend --target worker

# 2. Build the Next.js static export
cd frontend && npm run build && npm run export
# Output is in frontend/out/

# 3. Upload static export to S3 and invalidate CloudFront
aws s3 sync frontend/out/ s3://<your-bucket-name>/ --delete
aws cloudfront create-invalidation \
  --distribution-id <your-distribution-id> \
  --paths "/*"

# 4. Copy images or pull on the EC2 host
# (either push to ECR first, or scp the tarball)
docker save ecomm-api:latest | gzip | ssh ubuntu@<ec2-ip> docker load

# 5. On the EC2 host — decrypt env file (SOPS + age)
sops -d .env.enc > .env

# 6. Run migrations (one-shot container)
docker compose run --rm api goose -dir /app/migrations postgres \
  "${DATABASE_URL}" up

# 7. Restart services
docker compose up -d --force-recreate api worker

# 8. Smoke test
curl -f https://<api-domain>/healthz
curl -f https://<storefront-domain>/
```

If something goes wrong during deployment, roll back by restarting the previous image tag:

```bash
docker compose up -d --force-recreate api worker  # with old image tag pinned in compose.yml
```

---

*Guide generated: 2026-03-30. Maintained alongside `docs/01_prompts/06_dev_setup.prompt.md`.*
