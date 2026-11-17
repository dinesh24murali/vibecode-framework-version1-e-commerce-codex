# Changelog

All notable changes to this project will be documented here.

**AI tools: update this file on every task completion.**
Format: `- [YYYY-MM-DD] <task-name>: <one-line summary>`

---

## Unreleased

- [2026-11-17] task-020-contact-us: feat: POST /contact (stores enquiry, SMTP notification) and Contact Us page with map, social links and validated form
- [2026-11-17] task-019-static-pages: feat: static Wall of Fame, About Us, Terms and Privacy pages with unique titles
- [2026-11-17] task-018-footer: feat: shop footer with Terms, About, Wall of Fame links and env-configured social icons
- [2026-11-17] task-017-navbar-mega-menus: feat: shop navbar with Book Store and Printing Services mega-menus, search, cart badge and mobile accordion drawer
- [2026-11-17] task-016-product-detail-page: feat: /products/[slug] SSG page with image gallery, hover/pinch zoom, Add to Cart toast and per-product metadata
- [2026-11-17] task-015-product-list-page: feat: /products page with filter sidebar, search, sort, pagination, URL-synced filterSlice and mobile filter toggle
- [2026-11-17] task-014-landing-page: feat: GET /landing, carousel_slides table, static landing page with embla hero carousel and product grid (build-time data with client fallback)
- [2026-11-17] task-013-s3-presigned-urls: feat: admin image presign endpoint (S3 adapter, role middleware, env config, 503 when unconfigured)
- [2026-11-17] task-012-product-category-public-api: feat: public GET /products (filters, search, sort, pagination, filter facets), /products/{id}, /categories, /categories/{id}; migration for spec columns; real-DB filter tests
- [2026-11-17] task-011-attributes-system-api: feat: public GET /attributes and /attributes/{id}/values with contract and repository tests
- [2026-11-17] task-010-account-management: feat: account profile/password/email-change, address CRUD, order list/get/cancel endpoints and 3-tab account page; fix: removed debug prints that panicked on successful registration
- [2026-11-17] task-009-auth-pages-signup-login-facebook: feat: signup (email, OTP, password), login with ?redirect=, Facebook login button, Playwright auth e2e with mocked API, Playwright config
- [2026-11-17] task-040-background-worker-scaffold: feat: cmd/worker Redis Streams consumer (XREADGROUP, XACK, exponential backoff, dead-letter, idempotency keys), handler registry with Phase 2 stubs, worker compose service
- [2026-11-17] task-007-openapi-client-codegen-frontend: feat: orval-generated API client (lib/api/generated), base client with Bearer injection, cart slice wired to generated cart API, no-raw-fetch check in make verify
- [2026-11-17] task-008-zustand-store-auth-cart-slices: feat: Zustand auth + cart slices, SSG-safe persist (user profile only), useHydrateStore, injectable CartApi for syncWithServer
- [2026-11-17] task-000-docker-toolchain: chore: Docker-only toolchain (docker-compose.tools.yml, containerised Makefile targets and verify scripts), Postgres 18 volume fix, lint cleanup
- [2026-04-03] task-006-jwt-auth-system-backend: feat: full auth pipeline — Ed25519 JWT, Argon2id passwords, OTP flows (register + password-reset), Facebook OAuth, rotating refresh tokens, AuthMiddleware, 10 endpoints wired, ADR-0005, memory/auth.md, contract tests
- [2026-04-03] task-005-database-schema-full-migrations: chore: 20 goose migrations for full production schema (users, products, cart, orders, payments, coupons, consents, audit, addresses, contact_queries); 20 sqlc query files generated; Makefile auto-loads .env; Docker ports remapped (postgres 5433, redis 6380)
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
