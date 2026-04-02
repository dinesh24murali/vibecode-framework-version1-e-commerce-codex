# Changelog

All notable changes to this project will be documented here.

**AI tools: update this file on every task completion.**
Format: `- [YYYY-MM-DD] <task-name>: <one-line summary>`

---

## Unreleased

<!-- AI: add entries here as tasks are completed -->
- [2026-04-01] task-004-database-redis-setup-migrations: chore: pgxpool (MaxConns:20), goose migrations with embed, go-redis/v9, sqlc.yaml, migrate CLI, ADR-0004, memory/infra.md
- [2026-03-31] task-003-frontend-scaffold-nextjs: chore: NextJS 16.2.1 static export scaffold with App Router, Tailwind, shadcn/ui Button, Zustand v5 SSG-safe store, four route-group shell layouts, and ADR-0003
- [2026-03-31] task-002-backend-scaffold-go-gin: chore: Go 1.26/Gin backend scaffold with health endpoint, config validation, Builder/Adapter patterns, and internal package layout
- [2026-03-31] task-001-monorepo-scaffold-local-dev: chore: added docker-compose.yml, .env.example, backend/frontend Dockerfiles, and updated Makefile with all dev targets
- [2026-03-26] add-forgot-password-flow: docs: added 3-step forgot-password OTP flow to API spec and functional spec (UF-015, FS-006 expansion, notifications table, screen inventory)
- [2026-03-26] add-public-filter-options-api: docs: added GET /api/v1/attributes and GET /api/v1/attributes/{attrDefId}/values public endpoints plus AttributeWithValues schema to OpenAPI spec
- [2026-03-24] 004-correct-tech-architecture: docs: corrected the tech architecture schema, pagination, and frontend hosting model
- [2026-03-23] 001-generate-prd: docs: added the v1 product requirements document and product scope memory
- [2026-03-23] 002-generate-functional-spec: docs: added the functional specification and updated resolved product rules

---

## Legend
- `feat:` New feature
- `fix:` Bug fix
- `refactor:` Code change with no functional difference
- `docs:` Documentation only
- `test:` Test additions or changes
- `chore:` Build, config, or tooling changes
