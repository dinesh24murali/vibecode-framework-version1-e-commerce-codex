.PHONY: dev build test seed migrate migrate-up migrate-down migrate-status \
        sqlc-gen codegen verify adr task scaffold check-env init help

# Load .env if present so DATABASE_URL and friends are available to migrate targets.
ifneq (,$(wildcard .env))
  include .env
  export
endif

# ─── Local Development ────────────────────────────────────────────────────────

## dev: Start backend and frontend dev servers concurrently
dev:
	@echo "Starting backend and frontend dev servers..."
	@trap 'kill 0' INT; \
	  (cd backend && go run ./cmd/api) & \
	  (cd frontend && npm run dev) & \
	  wait

## build: Build all Docker images via Docker Compose
build:
	docker compose build

# ─── Testing ─────────────────────────────────────────────────────────────────

## test: Run all backend unit and integration tests
test:
	cd backend && go test ./...

## verify: Run all verification checks (e2e + DOM)
verify:
	@echo "Running verification..."
	@bash verify/scripts/run-e2e.sh
	@npx ts-node verify/scripts/check-dom.ts

# ─── Database ─────────────────────────────────────────────────────────────────

## migrate: Alias for migrate-up
migrate: migrate-up

## migrate-up: Apply all pending database migrations
migrate-up:
	cd backend && go run ./cmd/migrate up

## migrate-down: Roll back the most recent database migration
migrate-down:
	cd backend && go run ./cmd/migrate down

## migrate-status: Show current migration state
migrate-status:
	cd backend && go run ./cmd/migrate status

## seed: Seed the database with development fixture data
seed:
	cd backend && go run ./cmd/seed

## sqlc-gen: Regenerate Go database code from SQL query files
sqlc-gen:
	cd backend && sqlc generate

# ─── Code Generation ──────────────────────────────────────────────────────────

## codegen: Regenerate the frontend TypeScript API client from the OpenAPI spec
codegen:
	cd frontend && npx orval --config orval.config.ts

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
	 echo "Created adr/$$PADDED-$(SLUG).md"

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
