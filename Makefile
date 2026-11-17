.PHONY: dev build test seed migrate migrate-up migrate-down migrate-status \
        sqlc-gen codegen verify adr task scaffold check-env init help server-dev \
        go-build go-test go-lint fe-build fe-lint fe-install tools-build

# Load .env if present so DATABASE_URL and friends are available to migrate targets.
ifneq (,$(wildcard .env))
  include .env
  export
endif

# All Go/Node/sqlc/goose/playwright commands run in containers (see docker-compose.tools.yml).
# Git Bash on Windows mangles absolute container paths, so disable path conversion.
export MSYS_NO_PATHCONV=1
DC    := docker compose -f docker-compose.yml -f docker-compose.tools.yml
GO    := $(DC) run --rm go
NODE  := $(DC) run --rm node
E2E   := $(DC) run --rm e2e

# ─── Local Development ────────────────────────────────────────────────────────

## tools-build: Build the Go tools image (go, sqlc, goose, golangci-lint)
tools-build:
	$(DC) build go

## dev: Start postgres, redis, mailhog, api and frontend via Docker Compose
dev:
	docker compose up -d postgres redis mailhog api frontend

## build: Build all Docker images via Docker Compose
build:
	docker compose build

server-dev:
	docker compose up -d postgres redis mailhog api

# ─── Go (containerised) ──────────────────────────────────────────────────────

## go-build: go build ./... in the go container
go-build:
	$(GO) go build ./...

## go-test: go test ./... in the go container
go-test:
	$(GO) go test ./...

## go-lint: golangci-lint run in the go container
go-lint:
	$(GO) golangci-lint run

# ─── Frontend (containerised) ────────────────────────────────────────────────

## fe-install: npm ci in the node container (node_modules is a named volume)
fe-install:
	$(NODE) npm ci

## fe-build: Static export build of the frontend
fe-build: fe-install
	$(NODE) npm run build

## fe-lint: Type-check the frontend (next lint was removed in Next 16)
fe-lint: fe-install
	$(NODE) npx tsc --noEmit

# ─── Testing ─────────────────────────────────────────────────────────────────

## test: Run all backend unit and integration tests
test: go-test

## verify: Build the static export, serve it, run e2e + DOM checks (all in containers)
verify: fe-build
	@echo "Running verification..."
	@$(DC) up -d --force-recreate static
	@$(E2E) bash verify/scripts/run-e2e.sh
	@$(E2E) bash verify/scripts/run-check-dom.sh

# ─── Database ─────────────────────────────────────────────────────────────────

## migrate: Alias for migrate-up
migrate: migrate-up

## migrate-up: Apply all pending database migrations
migrate-up:
	$(GO) go run ./cmd/migrate up

## migrate-down: Roll back the most recent database migration
migrate-down:
	$(GO) go run ./cmd/migrate down

## migrate-status: Show current migration state
migrate-status:
	$(GO) go run ./cmd/migrate status

## seed: Seed the database with development fixture data
seed:
	$(GO) go run ./cmd/seed

## sqlc-gen: Regenerate Go database code from SQL query files
sqlc-gen:
	$(GO) sqlc generate

# ─── Code Generation ──────────────────────────────────────────────────────────

## codegen: Regenerate the frontend TypeScript API client from the OpenAPI spec
codegen:
	$(NODE) npx orval --config orval.config.ts

# ─── Project Scaffolding ──────────────────────────────────────────────────────

## scaffold: Create the initial project structure (run once after cloning)
scaffold:
	@echo "Scaffolding e-commerce-site..."
	@mkdir -p docs/02_outputs
	@mkdir -p tasks/active tasks/done
	@mkdir -p tests/e2e tests/contract
	@mkdir -p verify/scripts
	@mkdir -p backend frontend
	@echo "Done. Next: fill in docs/00_intake/intake_questionnaire.md"

## adr: Create a new ADR (usage: make adr SLUG=my-decision)
adr:
	@if [ -z "$(SLUG)" ]; then echo "Usage: make adr SLUG=my-decision"; exit 1; fi
	@N=$$(ls adr/*.md 2>/dev/null | grep -v template | grep -v AGENTS | wc -l); \
	 N=$$((N + 1)); \
	 PADDED=$$(printf "%04d" $$N); \
	 cp adr/template.md "adr/$$PADDED-$(SLUG).md"; \
	 sed -i -e "s/^# ADR-NNNN:/# ADR-$$PADDED:/" \
	        -e "s/^\*\*Status:\*\*.*/**Status:** Accepted/" \
	        -e "s/^\*\*Date:\*\* YYYY-MM-DD/**Date:** $$(date +%Y-%m-%d)/" \
	        "adr/$$PADDED-$(SLUG).md"; \
	 echo "Created adr/$$PADDED-$(SLUG).md (auto-accepted)"

## task: Create a new task from the template (usage: make task NAME=my-task)
task:
	@if [ -z "$(NAME)" ]; then echo "Usage: make task NAME=my-task"; exit 1; fi
	@cp tasks/_template.md "tasks/active/$(NAME).md"; \
	 echo "Created tasks/active/$(NAME).md"
	@echo "Next: ask your AI tool to write tasks/active/$(NAME).scratch.md before coding"

## check-env: Validate .env has all required vars from .env.example
check-env:
	@bash verify/scripts/check-env.sh

## init: List all unreplaced [[PLACEHOLDER]] tokens in the repo (run after cloning)
init:
	@echo "Scanning for unreplaced [[PLACEHOLDER]] tokens..."
	@echo ""
	@git ls-files | xargs grep -rn '\[\[' --include="*.md" --include="*.yml" --include="*.yaml" --include="*.sh" --include="*.ts" --include="*.json" 2>/dev/null | grep -v Binary || echo "✓ No unreplaced placeholders found"
	@echo ""
	@echo "Replace each token, then re-run 'make init' to verify."

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
