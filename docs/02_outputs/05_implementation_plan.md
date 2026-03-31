# Phased Implementation Plan — e-commerce-site (Book Store)

**Generated:** 2026-03-30
**Stack:** Go 1.26 / Gin · NextJS 16 SSG · PostgreSQL · JWT · AWS

---

> **Contract rule (AGENTS.md §3.7 + backend/AGENTS.md):**
> All endpoint paths and HTTP methods in this document are sourced from
> `docs/02_outputs/04_api_spec.yaml`. The spec is the authoritative contract.
> If any path below conflicts with the spec, the spec wins.
>
> **Password-reset flow** (`POST /api/v1/auth/password-reset/*`) is already
> implemented in `tasks/done/add-forgot-password-flow.md`. TASK-006 must
> acknowledge and wire those three endpoints but does not need to re-design them.
>
> **Per-task rules from AGENTS.md:**
> - Every backend task's acceptance criteria must include a contract test in `tests/contract/`
> - Every frontend task's acceptance criteria must include `make verify` passing
> - Every task file must list relevant `memory/*.md` files in its Dependencies

---

## Phase 0 — Project Setup

> Tasks are **sequential**. Nothing in Phase 1 can begin until TASK-004 is done.

---

#### TASK-001: Monorepo scaffold & local dev environment
- **Phase:** Phase 0 — Project Setup
- **Depends on:** none
- **Files:**
  - `Makefile`
  - `.gitignore`
  - `.env.example`
  - `docker-compose.yml`
  - `README.md`
- **Description:** Create the top-level monorepo with `backend/` and `frontend/` directories. Write a `docker-compose.yml` with all local dev services: `postgres` (port 5432), `redis` (port 6379), `api`, `worker`, `frontend`, and `mailhog` (ports 1025/8025). Add a root `Makefile` with targets: `dev`, `test`, `build`, `seed`, `migrate`, `migrate-status`, `sqlc-gen`, `codegen`, `verify`, `adr`, `task`.
- **Acceptance criteria:**
  - [ ] `docker-compose up` starts Postgres, Redis, MailHog, api, worker, and frontend with no errors
  - [ ] `make dev` starts both backend and frontend dev servers concurrently
  - [ ] `.env.example` documents every required environment variable with a one-line description (including `REDIS_URL`, `DATABASE_URL`, `MAILHOG_HOST`, `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET`)
  - [ ] `.gitignore` excludes `.env`, build artefacts, `vendor/`, `node_modules/`, `out/`

---

#### TASK-002: Backend scaffold — Go 1.26 / Gin
- **Phase:** Phase 0 — Project Setup
- **Depends on:** TASK-001
- **Files:**
  - `backend/main.go`
  - `backend/go.mod`
  - `backend/go.sum`
  - `backend/internal/config/config.go`
  - `backend/internal/server/server.go`
  - `backend/internal/server/router.go`
  - `backend/Makefile`
- **Description:** Initialise the Go module and install Gin. Wire a minimal HTTP server with a `/health` endpoint. Apply the Builder pattern for server construction and the Adapter pattern for config loading from environment variables. Establish the `internal/` package layout: `handlers/`, `services/`, `repositories/`, `models/`, `middleware/`, `db/`, `config/`.
- **Acceptance criteria:**
  - [ ] `go build ./...` succeeds with zero errors
  - [ ] `GET /health` returns `200 OK` with `{"status":"ok"}`
  - [ ] Config struct validates all required env vars at startup; server exits with a clear message if any are missing
  - [ ] `golangci-lint run` passes with zero warnings

> **ADR needed:** Go backend package layout and design-pattern conventions (Repository / Service / Handler layering, interface placement)

---

#### TASK-003: Frontend scaffold — NextJS 16 App Router
- **Phase:** Phase 0 — Project Setup
- **Depends on:** TASK-001
- **Files:**
  - `frontend/package.json`
  - `frontend/tsconfig.json`
  - `frontend/next.config.ts`
  - `frontend/tailwind.config.ts`
  - `frontend/app/layout.tsx`
  - `frontend/app/page.tsx`
  - `frontend/lib/store/index.ts`
  - `frontend/lib/store/slices/` *(empty directory with `.gitkeep`)*
- **Description:** Bootstrap NextJS 16 with App Router, TypeScript, Tailwind CSS, and shadcn/ui. Configure `output: 'export'` for static generation. Install Zustand v5 and create the slice-pattern store skeleton. Create route-group shell directories: `app/(shop)/`, `app/(auth)/`, `app/(account)/`, `app/(admin)/`. No pages yet — just empty `layout.tsx` files.
- **Acceptance criteria:**
  - [ ] `next build` produces a static export in `out/` with no errors
  - [ ] All four route-group directories exist with shell layouts
  - [ ] Zustand store initialises without errors in SSG context (no `localStorage` access during SSR)
  - [ ] `make verify` passes (baseline — no frontend code yet, just scaffold)

> **ADR needed:** Zustand SSR/SSG hydration strategy — how to prevent `localStorage` access during static generation

---

#### TASK-004: Database + Redis setup, migrations & query generation tooling
- **Phase:** Phase 0 — Project Setup
- **Depends on:** TASK-001
- **Files:**
  - `backend/migrations/` *(directory — goose migration files live here)*
  - `backend/internal/db/db.go`
  - `backend/internal/db/migrate.go`
  - `backend/internal/cache/redis.go`
  - `backend/db/queries/` *(directory — raw `.sql` query files consumed by sqlc)*
  - `backend/db/sqlc/` *(directory — sqlc-generated Go code; do not hand-edit)*
  - `sqlc.yaml`
- **Description:** Set up three foundational infrastructure layers:
  1. **Postgres (pgxpool):** Create `db.go` which opens a `pgxpool.Pool` (pgx v5) connection from `DATABASE_URL`. Configure pool limits: `MaxConns: 20`, `MinConns: 4`, `MaxConnLifetime: 30m`, `MaxConnIdleTime: 5m`. Server must refuse to start if the DB connection cannot be established.
  2. **Migrations (goose):** Integrate `pressly/goose` for SQL up/down migrations. Migration files in `backend/migrations/` must include `-- +goose Up` and `-- +goose Down` annotation comments. Create `migrate.go` which embeds the migrations directory and exposes a `RunMigrations(*pgxpool.Pool)` function. Migrations are run via a dedicated `make migrate-up` target — **not** auto-run on server startup. Add `make migrate-down` and `make migrate-status` targets.
  3. **Redis:** Create `cache/redis.go` which opens a `redis.Client` (go-redis/v9) from `REDIS_URL`. Used by auth (refresh token state), rate limiting, idempotency keys, and the job queue. Server must refuse to start if Redis is unreachable.
  4. **Query generation (sqlc):** Add `sqlc.yaml` pointing to `backend/migrations/` as the schema source and `backend/db/queries/` as the queries source. Configure the output package to `backend/db/sqlc/`. Add `make sqlc-gen` target. Repository implementations (`internal/repositories/`) must import from `backend/db/sqlc/` — never write raw SQL strings in Go code.
- **Acceptance criteria:**
  - [ ] `make migrate-up` runs all migrations against the local Postgres instance
  - [ ] `make migrate-down` rolls back the most recent migration
  - [ ] `make migrate-status` prints the current migration state
  - [ ] `make sqlc-gen` runs `sqlc generate` without errors (schema must exist first)
  - [ ] `sqlc.yaml` is committed; `backend/db/sqlc/` generated files are committed (or excluded via `.gitignore` with a note)
  - [ ] `pgxpool` connection honours the configured pool limits
  - [ ] Redis client connects and pings successfully on startup
  - [ ] `go test ./internal/db/...` and `go test ./internal/cache/...` verify connection state
  - [ ] Server logs an error and exits non-zero if `DATABASE_URL` or `REDIS_URL` is unreachable

---

#### TASK-040: Background worker scaffold
- **Phase:** Phase 0 — Project Setup
- **Depends on:** TASK-004
- **Files:**
  - `backend/cmd/worker/main.go`
  - `backend/internal/worker/worker.go`
  - `backend/internal/worker/dispatcher.go`
  - `backend/internal/worker/handlers/` *(empty directory with `.gitkeep`)*
- **Description:** Scaffold the Go background worker as a separate binary (`cmd/worker`). The worker connects to the same Postgres and Redis instances as the API. It consumes jobs from Redis Streams (`XREADGROUP`), dispatches to registered handlers by job type, and handles acknowledgement (`XACK`) on success or re-queues with backoff on failure. Job handler types to register in Phase 2: `order_confirmation_email`, `payment_reconciliation`, `stock_release`, `audit_fan_out`. The worker must be idempotent — replaying a job must not create duplicate side effects. Wire the worker as a separate service in `docker-compose.yml`.
- **Acceptance criteria:**
  - [ ] `go build ./cmd/worker` succeeds with zero errors
  - [ ] Worker connects to Redis and Postgres on startup; exits non-zero if either is unreachable
  - [ ] Worker consumes from a Redis Stream and dispatches to registered handlers
  - [ ] Failed jobs are requeued with exponential backoff (max 3 retries before dead-letter)
  - [ ] Worker runs as a separate `worker` service in `docker-compose.yml`
  - [ ] `golangci-lint run` passes

---

## Phase 1 — Core Infrastructure

> TASK-005 and TASK-007/008 can start in **parallel** once Phase 0 is done.
> TASK-006 depends on TASK-005 (schema must exist first).

---

#### TASK-005: Database schema — full migrations
- **Phase:** Phase 1 — Core Infrastructure
- **Depends on:** TASK-004
- **Files:**
  - `backend/migrations/0001_create_users.sql`
  - `backend/migrations/0002_create_refresh_tokens.sql`
  - `backend/migrations/0003_create_consents.sql`
  - `backend/migrations/0004_create_categories.sql`
  - `backend/migrations/0005_create_product_attributes.sql`
  - `backend/migrations/0006_create_product_attribute_values.sql`
  - `backend/migrations/0007_create_products.sql`
  - `backend/migrations/0008_create_product_images.sql`
  - `backend/migrations/0009_create_product_attribute_assignments.sql`
  - `backend/migrations/0010_create_inventory_items.sql`
  - `backend/migrations/0011_create_cart.sql`
  - `backend/migrations/0012_create_orders.sql`
  - `backend/migrations/0013_create_order_items.sql`
  - `backend/migrations/0014_create_order_address.sql`
  - `backend/migrations/0015_create_payments.sql`
  - `backend/migrations/0016_create_audit_events.sql`
  - `backend/migrations/0017_create_coupons.sql`
  - `backend/migrations/0018_create_coupon_users.sql`
  - `backend/migrations/0019_create_addresses.sql`
  - `backend/migrations/0020_create_contact_queries.sql`
- **Description:** Write all SQL migration files for goose. Each file must contain `-- +goose Up` and `-- +goose Down` annotation blocks. Key schema rules:
  - `users`: `email` as `CITEXT` (unique, case-insensitive); `role` CHECK constrained to `admin | customer`
  - `refresh_tokens`: stores `token_hash`, `token_family_id`, expiry, `revoked_at` — supports rotating refresh token family invalidation
  - `consents`: DPDPA compliance — records `consent_type`, `policy_version`, `granted_at`, `withdrawn_at NULL`
  - `product_attributes` (definitions) + `product_attribute_values` + `product_attribute_assignments` (join): flexible filter model
  - `products`: `discount_percent` CHECK `> 0 AND <= 100`; add a generated `tsvector` column combining `title`, `author_name`, `description`, `isbn13` with a GIN index for full-text search
  - `inventory_items`: one row per product; available stock = `on_hand - reserved`; concurrent mutations require row-level locking
  - `order_address`: immutable snapshot at checkout — no FK to `addresses`; `orders` never reference the `addresses` table directly
  - `orders`: `status` as a Postgres enum `pending_payment | paid | packed | shipped | delivered | cancelled | payment_failed | refunded`
  - `payments`: stores Razorpay `provider_order_id`, `provider_payment_id`, `gateway_payload JSONB`
  - `audit_events`: append-only event log per order
  - `coupons`: `type` column `flat | percentage`; `coupon_users` unique on `(user_id, coupon_id)`
  - After all migrations are written, add corresponding `.sql` query files in `backend/db/queries/` and run `make sqlc-gen`
- **Acceptance criteria:**
  - [ ] All 20 migrations apply cleanly with `make migrate-up`
  - [ ] Each migration file has valid `-- +goose Up` / `-- +goose Down` blocks
  - [ ] All FK constraints are enforced and verified with `psql`
  - [ ] `products` has a GIN index on the generated `tsvector` search column
  - [ ] `product_attribute_assignments` join table links products to attribute values
  - [ ] `inventory_items` exists with `on_hand`, `reserved`, `reorder_threshold`
  - [ ] `order_address` has no FK to `addresses`
  - [ ] `orders` status uses the Postgres enum; invalid values are rejected at DB level
  - [ ] `refresh_tokens` has index on `(user_id, expires_at)`
  - [ ] `payments` table has `gateway_payload JSONB` column
  - [ ] `make sqlc-gen` succeeds after all migrations are applied

---

#### TASK-006: JWT auth system — backend
- **Phase:** Phase 1 — Core Infrastructure
- **Depends on:** TASK-005
- **Files:**
  - `backend/internal/auth/jwt.go`
  - `backend/internal/auth/middleware.go`
  - `backend/internal/auth/otp.go`
  - `backend/internal/auth/facebook.go`
  - `backend/internal/handlers/auth_handler.go`
  - `backend/internal/services/auth_service.go`
  - `backend/internal/repositories/user_repository.go`
  - `tests/contract/auth_test.go`
- **Description:** Implement the full auth pipeline using spec-canonical paths. Registration: `POST /api/v1/auth/register/init` (send OTP) → `POST /api/v1/auth/register/verify-otp` (return short-lived verification token) → `POST /api/v1/auth/register/complete` (create user + issue JWT). Login: `POST /api/v1/auth/login`. Facebook: `POST /api/v1/auth/login/facebook`. Token lifecycle: `POST /api/v1/auth/refresh` and `POST /api/v1/auth/logout`. Password-reset endpoints (`/auth/password-reset/*`) are already spec'd in the done task — wire them here using the same OTP service. Follow Repository → Service → Handler layering with interfaces at each boundary (Bridge pattern between repository interface and concrete Postgres implementation).
- **Acceptance criteria:**
  - [ ] `POST /api/v1/auth/register/init` sends OTP (logged to stdout in dev); returns `200` even for unknown email to prevent enumeration
  - [ ] `POST /api/v1/auth/register/verify-otp` returns a short-lived `otp_token`; expires after 10 minutes
  - [ ] `POST /api/v1/auth/register/complete` requires valid `otp_token`; creates user and returns `AuthTokenResponse`
  - [ ] `POST /api/v1/auth/login` returns `access_token` + `refresh_token`
  - [ ] `POST /api/v1/auth/login/facebook` exchanges a Facebook access token for a JWT
  - [ ] `POST /api/v1/auth/refresh` issues a new access token; invalidates old refresh token
  - [ ] `POST /api/v1/auth/logout` invalidates the refresh token
  - [ ] All three `/auth/password-reset/*` routes are wired and return correct responses
  - [ ] `AuthMiddleware` rejects requests without a valid JWT with `401`
  - [ ] Contract test in `tests/contract/auth_test.go` covers all 10 auth endpoints against the spec

> **ADR needed:** OTP delivery mechanism — stdout logging in dev, SMTP vs third-party (SendGrid / AWS SES) in production

---

#### TASK-007: OpenAPI client codegen — frontend
- **Phase:** Phase 1 — Core Infrastructure
- **Depends on:** TASK-003
- **Files:**
  - `frontend/lib/api/` *(generated — do not hand-edit)*
  - `frontend/lib/api/client.ts` *(base client with auth header injection)*
  - `Makefile` *(update `codegen` target)*
- **Description:** Use `orval` or `openapi-typescript-codegen` to generate a fully typed API client from `docs/02_outputs/04_api_spec.yaml`. Wire the base URL from `NEXT_PUBLIC_API_URL`. The client must inject `Authorization: Bearer <token>` from the Zustand auth store when a token is present. Add `make codegen` target to `Makefile`. Update `frontend/AGENTS.md` to document the actual codegen command used.
- **Acceptance criteria:**
  - [ ] `make codegen` regenerates `frontend/lib/api/` without errors
  - [ ] Generated types match all schemas in the API spec
  - [ ] Auth header is injected automatically for protected endpoints
  - [ ] No hand-written `fetch` call exists anywhere in the frontend (enforced by `make verify`)
  - [ ] `frontend/AGENTS.md` `[[CODEGEN_COMMAND]]` placeholder is replaced with the actual command

---

#### TASK-008: Zustand v5 store — auth & cart slices
- **Phase:** Phase 1 — Core Infrastructure
- **Depends on:** TASK-003
- **Files:**
  - `frontend/lib/store/slices/authSlice.ts`
  - `frontend/lib/store/slices/cartSlice.ts`
  - `frontend/lib/store/index.ts`
- **Description:** Implement the `authSlice` (user profile, access token, refresh token, login/logout actions) and `cartSlice` (line items, quantities, add/remove/clear actions) following the Zustand v5 slice pattern. Persist auth tokens to `localStorage` via the Zustand `persist` middleware. The store must hydrate safely in the SSG context (no `localStorage` access on the server side).
- **Acceptance criteria:**
  - [ ] `authSlice` hydrates from `localStorage` on first client render without hydration mismatch warnings
  - [ ] `cartSlice` accumulates line items client-side; exposes `syncWithServer()` action for API sync
  - [ ] Both slices are combined in the root store with correct TypeScript types
  - [ ] `make verify` passes with no SSR hydration errors in the console

---

## Phase 2 — Feature Implementation

> Sub-phases 2A–2F are independent of each other once Phase 1 is complete.
> Within each sub-phase, backend and frontend tasks run **in parallel** once
> the API contract for that sub-phase is confirmed in the spec.

---

### 2A — Authentication & Account Management

#### TASK-009: Auth pages — signup, login, Facebook
- **Phase:** Phase 2 — Feature Implementation (2A)
- **Depends on:** TASK-006, TASK-007, TASK-008
- **Files:**
  - `frontend/app/(auth)/login/page.tsx`
  - `frontend/app/(auth)/signup/page.tsx`
  - `frontend/app/(auth)/signup/verify-otp/page.tsx`
  - `frontend/components/auth/LoginForm.tsx`
  - `frontend/components/auth/SignupForm.tsx`
  - `frontend/components/auth/OtpInput.tsx`
  - `frontend/components/auth/FacebookLoginButton.tsx`
  - `tests/e2e/auth.spec.ts`
- **Description:** Build the 3-step signup flow: (1) email entry → calls `registerInit`, (2) 6-digit OTP input → calls `registerVerifyOtp`, (3) password creation → calls `registerComplete`. Build the login page (email/password + Facebook button). On success, store tokens in `authSlice` and redirect to the page the user originally requested via `?redirect=` query param. Read `memory/auth.md` before starting.
- **Acceptance criteria:**
  - [ ] Full 3-step signup flow completes end-to-end against the running backend
  - [ ] Invalid OTP shows an inline error message; does not clear the OTP input
  - [ ] Facebook login button initiates the OAuth popup and resolves to a logged-in session
  - [ ] Auth tokens stored in `authSlice` and persisted to `localStorage`
  - [ ] `?redirect=` is honoured after successful login
  - [ ] E2e test in `tests/e2e/auth.spec.ts` covers signup and login happy paths
  - [ ] `make verify` passes

---

#### TASK-010: Account management — API + 3-tab page
- **Phase:** Phase 2 — Feature Implementation (2A)
- **Depends on:** TASK-009
- **Files:**
  - `backend/internal/handlers/account_handler.go`
  - `backend/internal/services/account_service.go`
  - `backend/internal/repositories/address_repository.go`
  - `backend/internal/repositories/order_repository.go`
  - `frontend/app/(account)/account/page.tsx`
  - `frontend/components/account/ProfileTab.tsx`
  - `frontend/components/account/AddressTab.tsx`
  - `frontend/components/account/OrdersTab.tsx`
  - `tests/contract/account_test.go`
  - `tests/e2e/account.spec.ts`
- **Description:** Implement backend endpoints — profile: `GET /api/v1/account/profile`, `PUT /api/v1/account/profile`, `PUT /api/v1/account/profile/password`, `POST /api/v1/account/profile/email-change/request`, `POST /api/v1/account/profile/email-change/confirm`; addresses: `GET /api/v1/account/addresses`, `POST /api/v1/account/addresses`, `PUT /api/v1/account/addresses/{id}`, `DELETE /api/v1/account/addresses/{id}`; orders: `GET /api/v1/orders`, `GET /api/v1/orders/{id}`, `PUT /api/v1/orders/{id}/cancel`. Build the 3-tab account page. Email change must re-trigger the OTP flow via `registerInit`.
- **Acceptance criteria:**
  - [ ] Profile tab saves first/last name and phone number
  - [ ] Email change triggers a new OTP; user must verify before the change is committed
  - [ ] Address tab creates, edits, and deletes addresses; all 9 fields (including GST) are saved
  - [ ] Orders tab shows a paginated list (10/page) with status badges (`open`, `fulfilled`, `cancelled`)
  - [ ] "Cancel" on an `open` order calls `PUT /api/v1/orders/{id}/cancel`; button is hidden for non-`open` orders
  - [ ] Unauthenticated access to `/account` redirects to `/login?redirect=/account`
  - [ ] Contract tests cover all account and order endpoints
  - [ ] `make verify` passes

---

### 2B — Product Catalog

#### TASK-011: Attributes system API
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-005
- **Files:**
  - `backend/internal/handlers/attribute_handler.go`
  - `backend/internal/services/attribute_service.go`
  - `backend/internal/repositories/attribute_repository.go`
  - `backend/internal/models/attribute.go`
  - `tests/contract/attributes_test.go`
- **Description:** Implement the two public attribute endpoints: `GET /api/v1/attributes` (returns all attribute definitions — level, theme, size, etc.) and `GET /api/v1/attributes/{attrDefId}/values` (returns values for a given definition). These power the product filter sidebar. Use the Abstract Factory pattern to build filter predicate objects from incoming query params.
- **Acceptance criteria:**
  - [ ] `GET /api/v1/attributes` returns all attribute definitions with `id`, `name`, `slug`
  - [ ] `GET /api/v1/attributes/{attrDefId}/values` returns values for the given definition
  - [ ] Both endpoints are public (no auth required)
  - [ ] Contract test in `tests/contract/attributes_test.go` validates response shapes against the spec
  - [ ] Seed data from TASK-030 can be queried successfully through these endpoints

---

#### TASK-012: Product & category public API
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-011
- **Files:**
  - `backend/internal/handlers/product_handler.go`
  - `backend/internal/handlers/category_handler.go`
  - `backend/internal/services/product_service.go`
  - `backend/internal/repositories/product_repository.go`
  - `backend/internal/repositories/category_repository.go`
  - `backend/internal/models/product.go`
  - `backend/internal/models/category.go`
  - `tests/contract/products_test.go`
- **Description:** Implement `GET /api/v1/products` with query-param filtering by `category`, `attributeValueIds[]`, and `search` (full-text on name/description). Implement `GET /api/v1/products/{id}`, `GET /api/v1/categories`, and `GET /api/v1/categories/{id}`. Responses include pagination metadata. Soft-deleted products (`deleted_at IS NOT NULL`) are excluded. Use the Abstract Factory pattern for query predicate construction.
- **Acceptance criteria:**
  - [ ] Filter by `category`, `attributeValueIds[]` and `search` return correct results verified against seed data
  - [ ] Pagination metadata (`total`, `page`, `per_page`, `total_pages`) is present in all list responses
  - [ ] `GET /api/v1/categories` returns all 12 seeded categories
  - [ ] Soft-deleted products do not appear in any public endpoint
  - [ ] Contract tests cover all 5 product/category endpoints

---

#### TASK-013: S3 presigned URLs for product images
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-002
- **Files:**
  - `backend/internal/services/storage_service.go`
  - `backend/internal/handlers/upload_handler.go`
  - `tests/contract/upload_test.go`
- **Description:** Implement an S3 Adapter (`StorageService` interface with an `S3Adapter` concrete type) that issues presigned `PUT` upload URLs and presigned `GET` read URLs. Expose `POST /api/v1/admin/upload/image/presign` (admin-only). Product records store only the S3 object key; presigned read URLs are generated on the fly in product API responses.
- **Acceptance criteria:**
  - [ ] `POST /api/v1/admin/upload/image/presign` returns a presigned PUT URL valid for 15 minutes
  - [ ] Unauthenticated or non-admin requests return `403`
  - [ ] Product list and detail responses include presigned GET URLs for all images
  - [ ] S3 bucket name, region, and credentials are read from environment variables — no hardcoded values
  - [ ] Contract test validates the `ImagePresignResponse` shape

---

#### TASK-014: Landing page — backend + frontend
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-012, TASK-007
- **Files:**
  - `backend/internal/handlers/landing_handler.go`
  - `backend/internal/services/landing_service.go`
  - `frontend/app/(shop)/page.tsx`
  - `frontend/components/shop/HeroCarousel.tsx`
  - `frontend/components/shop/ProductCard.tsx`
  - `frontend/components/shop/ProductGrid.tsx`
  - `tests/contract/landing_test.go`
  - `tests/e2e/landing.spec.ts`
- **Description:** Implement `GET /api/v1/landing` which returns carousel slides and featured products in one call. Build the frontend landing page as a statically generated server component: it calls `GET /api/v1/landing` at build time. The hero carousel uses `embla-carousel-react` with auto-advance and touch/swipe support. Below the carousel, render a `ProductGrid` of featured product cards.
- **Acceptance criteria:**
  - [ ] `GET /api/v1/landing` returns `LandingPageResponse` matching the spec schema
  - [ ] Carousel renders with 3–5 slides, auto-advances, and is swipeable on mobile
  - [ ] Product cards show name, price, discount badge, and "Add to Cart" button
  - [ ] Page is included in the static export (`out/index.html` exists after `next build`)
  - [ ] Contract test validates the `LandingPageResponse` schema
  - [ ] E2e test verifies the landing page loads and the carousel is visible
  - [ ] `make verify` passes

---

#### TASK-015: Product list page with filters
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-012, TASK-011, TASK-007
- **Files:**
  - `frontend/app/(shop)/products/page.tsx`
  - `frontend/components/shop/FilterSidebar.tsx`
  - `frontend/components/shop/ProductListGrid.tsx`
  - `frontend/lib/store/slices/filterSlice.ts`
  - `tests/e2e/product-list.spec.ts`
- **Description:** Build the product list page. The left `FilterSidebar` fetches attribute definitions from `GET /api/v1/attributes` and renders checkbox groups per definition. Active filters are serialised into URL search params (`?category=...&attributeValueIds=...&search=...`). The `filterSlice` hydrates from `useSearchParams()` on mount and updates the URL on every filter change. Product grid updates via client-side API call — no full page reload.
- **Acceptance criteria:**
  - [ ] Active filters are reflected in the URL: `/products?category=Journals&attributeValueIds=3,7`
  - [ ] Navigating directly to a filtered URL pre-checks the correct sidebar filters on load
  - [ ] Filter changes update the product grid without a full page reload
  - [ ] Mobile: filter sidebar collapses behind a "Filters" drawer/sheet toggle
  - [ ] E2e test verifies filtering and URL sync
  - [ ] `make verify` passes

---

#### TASK-016: Product detail page
- **Phase:** Phase 2 — Feature Implementation (2B)
- **Depends on:** TASK-012, TASK-007
- **Files:**
  - `frontend/app/(shop)/products/[slug]/page.tsx`
  - `frontend/components/shop/ProductImageGallery.tsx`
  - `frontend/components/shop/ProductZoom.tsx`
  - `frontend/components/shop/AddToCartButton.tsx`
  - `tests/e2e/product-detail.spec.ts`
- **Description:** Build the product detail page using `generateStaticParams` to pre-render all product slugs at build time from `GET /api/v1/products`. Implement an image gallery with thumbnail strip and zoom-on-hover (desktop) / pinch-to-zoom (mobile) using `react-image-magnify` or equivalent. A single "Add to Cart" button updates the Zustand `cartSlice`.
- **Acceptance criteria:**
  - [ ] All product slugs are pre-rendered at build time (`generateStaticParams` fetches from API)
  - [ ] Image zoom works on desktop hover; pinch-to-zoom works on mobile
  - [ ] "Add to Cart" updates `cartSlice` and shows a toast notification
  - [ ] Page includes correct `<title>` and `<meta name="description">` from product data
  - [ ] E2e test covers load, image gallery interaction, and add-to-cart
  - [ ] `make verify` passes

---

### 2C — Navigation & Static Pages

#### TASK-017: Navbar with mega-menus
- **Phase:** Phase 2 — Feature Implementation (2C)
- **Depends on:** TASK-003
- **Files:**
  - `frontend/components/layout/Navbar.tsx`
  - `frontend/components/layout/BookStoreMegaMenu.tsx`
  - `frontend/components/layout/PrintingServicesMegaMenu.tsx`
  - `frontend/components/layout/SearchBar.tsx`
  - `frontend/lib/constants/navigation.ts`
  - `tests/e2e/navbar.spec.ts`
- **Description:** Build the responsive navbar. `BookStoreMegaMenu` (hover-activated) lists all 12 product categories and links to `/products?category=<name>`. `PrintingServicesMegaMenu` shows two columns: Products Offered (11 items) and Inhouse Services (13 items). `SearchBar` (top-right) submits on `Enter` to `/products?search=<query>`. Navbar links: Book Store, Wall of Fame, Contact Us, Printing Services, Printing Supplies. Cart icon shows item count from `cartSlice`.
- **Acceptance criteria:**
  - [ ] Hovering "Book Store" reveals all 12 categories; clicking one navigates to `/products?category=<name>`
  - [ ] Hovering "Printing Services" reveals both columns with all 24 sub-items
  - [ ] Search bar submits on `Enter` and navigates to `/products?search=<query>`
  - [ ] Cart icon badge reflects `cartSlice` item count
  - [ ] Mobile: mega-menus collapse into a hamburger drawer with accordion sections
  - [ ] E2e test covers category click and search submission
  - [ ] `make verify` passes

---

#### TASK-018: Footer component
- **Phase:** Phase 2 — Feature Implementation (2C)
- **Depends on:** TASK-003
- **Files:**
  - `frontend/components/layout/Footer.tsx`
- **Description:** Build the site footer with internal links (Terms & Conditions, About Us, Wall of Fame) and external social media icon links (Facebook, Instagram, WhatsApp). Use `next/link` for internal routes and `target="_blank" rel="noopener noreferrer"` for external links.
- **Acceptance criteria:**
  - [ ] All internal links navigate correctly
  - [ ] Social media icons open in a new tab
  - [ ] Footer is present on every page via the root `(shop)` layout
  - [ ] `make verify` passes

---

#### TASK-019: Static pages — Wall of Fame, About Us, T&C, Privacy Policy
- **Phase:** Phase 2 — Feature Implementation (2C)
- **Depends on:** TASK-003
- **Files:**
  - `frontend/app/(shop)/wall-of-fame/page.tsx`
  - `frontend/app/(shop)/about/page.tsx`
  - `frontend/app/(shop)/terms/page.tsx`
  - `frontend/app/(shop)/privacy/page.tsx`
- **Description:** Create four static server-component pages. Wall of Fame and About Us use rich placeholder content with dummy images. Terms & Conditions and Privacy Policy contain boilerplate legal sections with proper headings. All four are included in the static export.
- **Acceptance criteria:**
  - [ ] All four routes exist in the static export (`out/` contains each path)
  - [ ] Each page has a unique and descriptive `<title>` tag
  - [ ] No client-side JavaScript errors on any of the four pages
  - [ ] `make verify` passes

---

#### TASK-020: Contact Us — API + page
- **Phase:** Phase 2 — Feature Implementation (2C)
- **Depends on:** TASK-003, TASK-005
- **Files:**
  - `backend/internal/handlers/contact_handler.go`
  - `backend/internal/services/contact_service.go`
  - `frontend/app/(shop)/contact/page.tsx`
  - `frontend/components/shop/ContactForm.tsx`
  - `tests/contract/contact_test.go`
  - `tests/e2e/contact.spec.ts`
- **Description:** Implement `POST /api/v1/contact` which stores a contact query (name, email, phone, message) in the `contact_queries` table and returns `201`. Build the Contact Us page with address block, embedded Google Maps iframe, social media links (Facebook, Instagram, WhatsApp), and the contact form.
- **Acceptance criteria:**
  - [ ] Form submission stores the query in DB and the API returns `201`
  - [ ] Google Maps iframe is embedded with the correct address
  - [ ] Form shows a success message on `201` and an error message on API failure
  - [ ] All form fields have client-side validation (required, email format, phone format)
  - [ ] Contract test validates the `ContactQueryRequest` schema and `201` response
  - [ ] E2e test covers happy-path form submission
  - [ ] `make verify` passes

---

### 2D — Cart & Checkout

#### TASK-021: Cart API
- **Phase:** Phase 2 — Feature Implementation (2D)
- **Depends on:** TASK-006, TASK-005
- **Files:**
  - `backend/internal/handlers/cart_handler.go`
  - `backend/internal/services/cart_service.go`
  - `backend/internal/repositories/cart_repository.go`
  - `tests/contract/cart_test.go`
- **Description:** Implement all cart endpoints from the spec: `GET /api/v1/cart`, `DELETE /api/v1/cart` (clear cart), `POST /api/v1/cart/items` (add item), `PUT /api/v1/cart/items/{itemId}` (update quantity), `DELETE /api/v1/cart/items/{itemId}` (remove item). Cart is user-scoped and requires authentication. Adding an item that already exists increments its quantity.
- **Acceptance criteria:**
  - [ ] Unauthenticated requests return `401`
  - [ ] Adding an existing product increments its quantity rather than creating a duplicate line
  - [ ] `DELETE /api/v1/cart` clears all items; `GET /api/v1/cart` returns an empty cart after
  - [ ] Cart total is computed server-side including any per-product discounts
  - [ ] `PUT /api/v1/cart/items/{itemId}` with `qty: 0` is rejected with `400`
  - [ ] Contract tests cover all 5 cart endpoints

---

#### TASK-022: Cart page — frontend
- **Phase:** Phase 2 — Feature Implementation (2D)
- **Depends on:** TASK-021, TASK-007, TASK-008
- **Files:**
  - `frontend/app/(shop)/cart/page.tsx`
  - `frontend/components/shop/CartLineItem.tsx`
  - `frontend/components/shop/CartSummary.tsx`
  - `tests/e2e/cart.spec.ts`
- **Description:** Build the cart page. On mount, authenticated users sync `cartSlice` with the server via `GET /api/v1/cart`. Quantity stepper calls `PUT /api/v1/cart/items/{itemId}` and removes the item when quantity reaches zero. Cart summary shows subtotal, discount, and a "Proceed to Checkout" button. Empty cart shows an illustration and "Continue Shopping" link.
- **Acceptance criteria:**
  - [ ] Quantity changes are optimistically updated in `cartSlice` then confirmed by the API
  - [ ] Empty cart state is shown when the cart has no items
  - [ ] Cart persists across page refreshes for authenticated users
  - [ ] "Proceed to Checkout" is disabled / shows a login prompt for unauthenticated users
  - [ ] E2e test covers add-to-cart from product page, cart page view, and quantity update
  - [ ] `make verify` passes

---

#### TASK-023: Checkout API — calculate, coupon, place order
- **Phase:** Phase 2 — Feature Implementation (2D)
- **Depends on:** TASK-021
- **Files:**
  - `backend/internal/handlers/checkout_handler.go`
  - `backend/internal/services/checkout_service.go`
  - `backend/internal/services/coupon_service.go`
  - `backend/internal/repositories/coupon_repository.go`
  - `tests/contract/checkout_test.go`
- **Description:** Implement three checkout endpoints: `POST /api/v1/checkout/calculate` (returns itemised subtotal, tax, discount, grand total without creating an order), `POST /api/v1/checkout/coupons/validate` (validates a coupon code and returns the discount), and `POST /api/v1/checkout/orders` (creates the order, decrements stock atomically, returns the order record). Use the Builder pattern for order assembly. Apply 18% GST on taxable items. Support flat (`SUPERHIT`: ₹150 off) and percentage (`FIRSTTIME`: 10% off) coupons.
- **Acceptance criteria:**
  - [ ] `POST /api/v1/checkout/calculate` returns a `CheckoutSummary` with correct tax line
  - [ ] `POST /api/v1/checkout/coupons/validate` returns `200` with discount amount for valid codes; `404` for unknown codes
  - [ ] `POST /api/v1/checkout/orders` creates an order and atomically decrements product stock
  - [ ] Placing an order with an empty cart returns `400`
  - [ ] Coupon deduction is applied before tax calculation
  - [ ] Contract tests cover all three endpoints with both coupon types and edge cases

---

#### TASK-039: Razorpay payment integration — backend
- **Phase:** Phase 2 — Feature Implementation (2D)
- **Depends on:** TASK-023
- **Files:**
  - `backend/internal/handlers/payment_handler.go`
  - `backend/internal/services/payment_service.go`
  - `backend/internal/repositories/payment_repository.go`
  - `tests/contract/payments_test.go`
- **Description:** Implement the full Razorpay payment lifecycle. After `POST /api/v1/checkout/orders` creates an order, the frontend must initiate payment via Razorpay JS SDK. Backend responsibilities:
  1. `POST /api/v1/payments/initiate` — creates a Razorpay order (via Razorpay API v1), stores `provider_order_id` in the `payments` table, returns the Razorpay order ID and key to the frontend. Requires `Idempotency-Key` header.
  2. `POST /api/v1/payments/webhooks/razorpay` — verifies the Razorpay webhook HMAC signature (`X-Razorpay-Signature`); on `payment.captured`: updates `payments.status`, transitions `orders.status` to `paid`, decrements `inventory_items.reserved`, publishes an order-confirmation job to the Redis Streams queue. Return `200` immediately; all side effects run asynchronously via worker. On `payment.failed`: transitions order to `payment_failed`.
  - `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET`, and `RAZORPAY_WEBHOOK_SECRET` are read from environment variables.
- **Acceptance criteria:**
  - [ ] `POST /api/v1/payments/initiate` returns a valid Razorpay order ID
  - [ ] Webhook with invalid signature returns `400` immediately
  - [ ] Webhook `payment.captured` transitions order status to `paid` and updates `payments` table
  - [ ] Webhook `payment.failed` transitions order status to `payment_failed`
  - [ ] Payment handler is idempotent — replaying the same `Idempotency-Key` returns the cached response
  - [ ] No secrets are hardcoded; all keys read from env vars
  - [ ] Contract tests cover both payment endpoints and webhook validation

---

#### TASK-024: Checkout page — frontend
- **Phase:** Phase 2 — Feature Implementation (2D)
- **Depends on:** TASK-023, TASK-022
- **Files:**
  - `frontend/app/(shop)/checkout/page.tsx`
  - `frontend/components/shop/CheckoutSummary.tsx`
  - `frontend/components/shop/CouponInput.tsx`
  - `frontend/components/shop/AddressSelector.tsx`
  - `tests/e2e/checkout.spec.ts`
- **Description:** Build the checkout page. Guard: redirect to `/cart` if `cartSlice` is empty; redirect to `/login?redirect=/checkout` if unauthenticated. The page calls `POST /api/v1/checkout/calculate` on load for the initial price breakdown. `CouponInput` calls `POST /api/v1/checkout/coupons/validate` on blur and updates the summary live. `AddressSelector` lets the user pick from their saved addresses. "Place Order" calls `POST /api/v1/checkout/orders`.
- **Acceptance criteria:**
  - [ ] Empty-cart redirect works when navigating directly to `/checkout`
  - [ ] Unauthenticated redirect works
  - [ ] Coupon field shows live valid/invalid feedback on blur; total updates when a valid coupon is applied
  - [ ] Successful order placement navigates to `/account?tab=orders`
  - [ ] E2e test covers the full checkout happy path (add to cart → checkout → coupon → place order)
  - [ ] `make verify` passes

---

### 2E — Admin Panel

#### TASK-025: Admin auth — middleware + login page
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-006
- **Files:**
  - `backend/internal/auth/admin_middleware.go`
  - `frontend/app/(admin)/admin/login/page.tsx`
  - `frontend/components/admin/AdminLoginForm.tsx`
  - `frontend/lib/store/slices/adminAuthSlice.ts`
  - `tests/contract/admin_auth_test.go`
  - `tests/e2e/admin-login.spec.ts`
- **Description:** Add `AdminMiddleware` that verifies the `role: admin` JWT claim and returns `403` for non-admin tokens. Implement `POST /api/v1/admin/auth/login` (email/password only, no OTP). Build the admin login page at `/admin/login` with its own layout (no public navbar/footer). Admin session stored in `adminAuthSlice`. All `/admin/*` pages redirect to `/admin/login` when unauthenticated.
- **Acceptance criteria:**
  - [ ] Admin login returns a JWT with `role: admin` claim
  - [ ] Non-admin JWTs receive `403` on all `/api/v1/admin/*` routes
  - [ ] Admin login page uses a separate layout with no public site chrome
  - [ ] Unauthenticated access to any `/admin/*` page redirects to `/admin/login`
  - [ ] Contract test validates the `POST /api/v1/admin/auth/login` response
  - [ ] E2e test covers admin login happy path and redirect behaviour
  - [ ] `make verify` passes

---

#### TASK-026: Admin — attributes CRUD
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-025, TASK-011
- **Files:**
  - `backend/internal/handlers/admin_attribute_handler.go`
  - `frontend/app/(admin)/admin/attributes/page.tsx`
  - `frontend/components/admin/AttributeTable.tsx`
  - `frontend/components/admin/AttributeFormModal.tsx`
  - `frontend/components/admin/AttributeValueFormModal.tsx`
  - `tests/contract/admin_attributes_test.go`
- **Description:** Implement the full admin attributes CRUD as defined in the spec. Attribute definitions: `GET /api/v1/admin/attributes`, `POST /api/v1/admin/attributes`, `GET /api/v1/admin/attributes/{attrDefId}`, `PUT /api/v1/admin/attributes/{attrDefId}`, `DELETE /api/v1/admin/attributes/{attrDefId}`. Attribute values: `GET /api/v1/admin/attributes/{attrDefId}/values`, `POST /api/v1/admin/attributes/{attrDefId}/values`, `PUT /api/v1/admin/attributes/{attrDefId}/values/{valueId}`, `DELETE /api/v1/admin/attributes/{attrDefId}/values/{valueId}`. Build the admin attributes page using the modal pattern (no dynamic `[id]` route — incompatible with `output: 'export'`). This manages level, theme, size, and any future filter dimensions.
- **Acceptance criteria:**
  - [ ] Admin can create, rename, and delete attribute definitions (e.g. "Level", "Theme", "Size")
  - [ ] Admin can add, rename, and delete values per definition (e.g. "Beginner", "Advanced")
  - [ ] Deleting an attribute value that is assigned to products returns `409 Conflict`
  - [ ] All CRUD operations work via modals; no `[id]` dynamic route exists
  - [ ] Contract tests cover all 9 attribute endpoints
  - [ ] `make verify` passes

---

#### TASK-027: Admin — product CRUD
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-025, TASK-012, TASK-013
- **Files:**
  - `backend/internal/handlers/admin_product_handler.go`
  - `frontend/app/(admin)/admin/products/page.tsx`
  - `frontend/components/admin/ProductTable.tsx`
  - `frontend/components/admin/ProductFormModal.tsx`
  - `tests/contract/admin_products_test.go`
- **Description:** Implement admin product endpoints: `GET /api/v1/admin/products`, `POST /api/v1/admin/products`, `GET /api/v1/admin/products/{id}`, `PUT /api/v1/admin/products/{id}`, `DELETE /api/v1/admin/products/{id}` (soft-delete), `GET /api/v1/admin/products/{id}/attributes`, `POST /api/v1/admin/products/{id}/attributes`, `DELETE /api/v1/admin/products/{id}/attributes/{attrId}`. Build the products page with a data table and search. "Edit" and "New Product" open `ProductFormModal` — a drawer that fetches `GET /api/v1/admin/products/{id}` on open and embeds the S3 image uploader (presigned PUT URL from TASK-013). No `[id]` dynamic route (SSG incompatibility).
- **Acceptance criteria:**
  - [ ] Admin can create, edit, soft-delete products from the list page via modal
  - [ ] "Edit" row action opens the modal pre-populated via `GET /api/v1/admin/products/{id}`
  - [ ] Image upload uses the presigned PUT URL; S3 key is saved to the product record
  - [ ] Attribute assignment (`POST /api/v1/admin/products/{id}/attributes`) links attribute values to the product
  - [ ] Attribute removal (`DELETE /api/v1/admin/products/{id}/attributes/{attrId}`) unlinks a single attribute value
  - [ ] Soft-deleted products are excluded from public endpoints
  - [ ] No `[id]` dynamic route segment exists under `app/(admin)/admin/products/`
  - [ ] Contract tests cover all 8 admin product endpoints
  - [ ] `make verify` passes

---

#### TASK-028: Admin — category CRUD
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-025, TASK-012
- **Files:**
  - `backend/internal/handlers/admin_category_handler.go`
  - `frontend/app/(admin)/admin/categories/page.tsx`
  - `frontend/components/admin/CategoryFormModal.tsx`
  - `tests/contract/admin_categories_test.go`
- **Description:** Implement admin category endpoints: `GET /api/v1/admin/categories`, `POST /api/v1/admin/categories`, `GET /api/v1/admin/categories/{id}`, `PUT /api/v1/admin/categories/{id}`, `DELETE /api/v1/admin/categories/{id}`. Build the categories page with a simple list and create/edit/delete actions via modal (no `[id]` dynamic route). Deleting a category with linked products returns `409 Conflict`.
- **Acceptance criteria:**
  - [ ] Admin can create, rename, and delete categories
  - [ ] Deleting a category with linked active products returns `409` with a clear message
  - [ ] Category changes are immediately reflected in the list without a page reload
  - [ ] No `[id]` dynamic route exists under `app/(admin)/admin/categories/`
  - [ ] Contract tests cover all 5 category endpoints
  - [ ] `make verify` passes

---

#### TASK-029: Admin — orders management
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-025, TASK-010
- **Files:**
  - `backend/internal/handlers/admin_order_handler.go`
  - `frontend/app/(admin)/admin/orders/page.tsx`
  - `frontend/components/admin/OrderTable.tsx`
  - `frontend/components/admin/OrderDetailModal.tsx`
  - `tests/contract/admin_orders_test.go`
- **Description:** Implement `GET /api/v1/admin/orders` (paginated, filterable by status), `GET /api/v1/admin/orders/{id}`, and `PUT /api/v1/admin/orders/{id}/status`. Build the admin orders page: data table with status filter tabs and a row action to open the order detail modal. Status can be changed to `fulfilled` or `cancelled` from the modal dropdown.
- **Acceptance criteria:**
  - [ ] Admin can filter orders by `open`, `fulfilled`, `cancelled` status tabs
  - [ ] Status can be changed to `fulfilled` or `cancelled` via the modal dropdown
  - [ ] Order detail modal shows customer name, delivery address, line items, and total
  - [ ] Table is paginated (20/page)
  - [ ] Contract tests cover all 3 admin order endpoints
  - [ ] `make verify` passes

---

#### TASK-041: Admin — coupon CRUD
- **Phase:** Phase 2 — Feature Implementation (2E)
- **Depends on:** TASK-025
- **Files:**
  - `backend/internal/handlers/admin_coupon_handler.go`
  - `backend/internal/services/coupon_service.go` *(extend from TASK-023)*
  - `backend/internal/repositories/coupon_repository.go` *(extend from TASK-023)*
  - `frontend/app/(admin)/admin/coupons/page.tsx`
  - `frontend/components/admin/CouponTable.tsx`
  - `frontend/components/admin/CouponFormModal.tsx`
  - `tests/contract/admin_coupons_test.go`
- **Description:** Implement admin coupon management endpoints: `GET /api/v1/admin/coupons`, `POST /api/v1/admin/coupons`, `GET /api/v1/admin/coupons/{id}`, `PUT /api/v1/admin/coupons/{id}`, `DELETE /api/v1/admin/coupons/{id}`. Each coupon has `code`, `type` (`flat | percentage`), `flat` amount or `percent` value, `is_enable` toggle, and `coupon_expiry`. Build the admin coupons page with a data table and create/edit/toggle-enable actions via modal (no `[id]` dynamic route — SSG constraint). Deleting a coupon that has been redeemed (`coupon_users` rows exist) returns `409 Conflict`.
- **Acceptance criteria:**
  - [ ] Admin can create, edit, enable/disable, and delete coupons
  - [ ] Coupon `type` field correctly gates whether `flat` or `percent` amount is used
  - [ ] Deleting a redeemed coupon returns `409` with a clear message
  - [ ] No `[id]` dynamic route exists under `app/(admin)/admin/coupons/`
  - [ ] Contract tests cover all 5 admin coupon endpoints
  - [ ] `make verify` passes

---

### 2F — Seed Scripts

#### TASK-030: Database seed scripts
- **Phase:** Phase 2 — Feature Implementation (2F)
- **Depends on:** TASK-005, TASK-006
- **Files:**
  - `backend/cmd/seed/main.go`
  - `backend/cmd/seed/seeds/coupons.go`
  - `backend/cmd/seed/seeds/admin.go`
  - `backend/cmd/seed/seeds/categories.go`
  - `backend/cmd/seed/seeds/attributes.go`
  - `backend/cmd/seed/seeds/products.go`
- **Description:** Write an idempotent Go seed binary (`cmd/seed`) that inserts: (1) coupons `SUPERHIT` (flat ₹150 off) and `FIRSTTIME` (10% off), (2) the default admin user (email + bcrypt-hashed password from `ADMIN_DEFAULT_PASSWORD` env var), (3) all 12 product categories, (4) attribute definitions for Level, Theme, and Size with representative values, (5) 5 dummy products across 3 of those categories, each with images, stock, and attribute assignments.
- **Acceptance criteria:**
  - [ ] `make seed` runs without errors on a freshly migrated DB
  - [ ] Re-running `make seed` does not create duplicate records (upsert logic)
  - [ ] Admin password is read from `ADMIN_DEFAULT_PASSWORD` env var — never hardcoded
  - [ ] All 5 dummy products have at least one image key, stock > 0, and attribute values assigned
  - [ ] `GET /api/v1/categories` returns all 12 categories after seeding
  - [ ] `GET /api/v1/attributes` returns Level, Theme, and Size definitions after seeding

---

## Phase 3 — Quality & Hardening

> All four tasks are **parallel** and can begin once Phase 2 is feature-complete.
> TASK-031 and TASK-032 are backend-only; TASK-034 is frontend-only.

---

#### TASK-031: Backend input validation, structured error handling & idempotency
- **Phase:** Phase 3 — Quality & Hardening
- **Depends on:** TASK-012, TASK-021, TASK-023
- **Files:**
  - `backend/internal/middleware/error_handler.go`
  - `backend/internal/middleware/request_id.go`
  - `backend/internal/middleware/idempotency.go`
  - `backend/internal/errors/errors.go`
- **Description:** Add three middleware layers:
  1. **Error handler:** Global Gin middleware that converts all domain errors to a consistent JSON envelope matching the spec's `ErrorResponse` schema: `{"error": {"code": "...", "message": "...", "fields": [...]}}`. Integrate `go-playground/validator` for request body validation on all POST/PUT/PATCH handlers.
  2. **Request ID:** Generates and propagates a UUID in `X-Request-ID` on every response.
  3. **Idempotency key:** `IdempotencyMiddleware` reads `Idempotency-Key` header on checkout and payment-initiation endpoints. Stores the key in Redis with a 24-hour TTL; replays the cached response if the same key is seen again. Returns `422` if the header is missing on required endpoints.
- **Acceptance criteria:**
  - [ ] All 4xx and 5xx responses conform to the `ErrorResponse` schema
  - [ ] Missing required fields return `422` with per-field `FieldError` objects
  - [ ] Unhandled panics are recovered and return `500` with the request ID
  - [ ] `X-Request-ID` header is present on every response
  - [ ] `POST /api/v1/checkout/orders` and payment endpoints require `Idempotency-Key`; missing key returns `422`
  - [ ] Replaying the same `Idempotency-Key` returns the cached response without re-processing
  - [ ] `go test ./internal/middleware/...` passes

---

#### TASK-032: Rate limiting & security headers
- **Phase:** Phase 3 — Quality & Hardening
- **Depends on:** TASK-006
- **Files:**
  - `backend/internal/middleware/rate_limit.go`
  - `backend/internal/middleware/security_headers.go`
- **Description:** Add Redis-backed per-IP rate limiting on auth endpoints using a sliding-window counter stored in Redis (go-redis/v9 + `INCR`/`EXPIRE`): 5 req/min on `/register/init` and `/password-reset/request`; 20 req/min on `/login`. Using Redis (not in-memory) ensures rate limits survive process restarts and work correctly if multiple API replicas are ever run. Add a security headers middleware: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Strict-Transport-Security`, `Content-Security-Policy`. CORS must restrict origins to the production frontend URL in `APP_ENV=production`.
- **Acceptance criteria:**
  - [ ] Exceeding OTP rate limit returns `429 Too Many Requests` with a `Retry-After` header
  - [ ] Rate limit counters are stored in Redis, not in-process memory
  - [ ] All responses include the four required security headers
  - [ ] CORS allows only the frontend origin in production; is open in development
  - [ ] `go test ./internal/middleware/...` covers rate-limit boundary cases

---

#### TASK-033: Backend test coverage
- **Phase:** Phase 3 — Quality & Hardening
- **Depends on:** TASK-006, TASK-012, TASK-023
- **Files:**
  - `backend/internal/services/*_test.go`
  - `backend/internal/repositories/*_test.go`
  - `tests/coverage-baseline.json`
- **Description:** Write unit tests for all service-layer functions (mock repositories via interfaces). Write integration tests for auth, product, cart, and checkout handlers using a real test Postgres instance (per `tests/AGENTS.md` — no DB mocking). Update `tests/coverage-baseline.json` with a minimum of 70% line coverage per package.
- **Acceptance criteria:**
  - [ ] `go test ./...` passes with zero failures
  - [ ] Line coverage ≥ 70% across all `internal/` packages (verified by `go test -cover`)
  - [ ] Auth OTP flow has a full integration test: `register/init` → `register/verify-otp` → `register/complete` → `login`
  - [ ] Checkout coupon logic has table-driven tests covering flat, percentage, invalid, and expired codes
  - [ ] `tests/coverage-baseline.json` is updated with thresholds for all new packages

---

#### TASK-034: Mobile responsiveness audit
- **Phase:** Phase 3 — Quality & Hardening
- **Depends on:** TASK-014, TASK-015, TASK-016, TASK-017, TASK-022, TASK-024
- **Files:**
  - `frontend/components/**/*.tsx` *(targeted fixes)*
  - `verify/dom-checks.md` *(update checklist)*
- **Description:** Run a full layout audit across all pages at 375px, 768px, and 1280px breakpoints using Playwright. Fix any horizontal overflow, touch-target size, or stacking issues found. Update `verify/dom-checks.md` with the mobile checklist items verified.
- **Acceptance criteria:**
  - [ ] No horizontal scroll at 375px on any page
  - [ ] All interactive elements meet a minimum 44×44px touch target
  - [ ] Navbar hamburger menu works correctly on both iOS Safari and Android Chrome viewports
  - [ ] Carousel is swipeable on touch device viewports
  - [ ] `make verify` passes with updated e2e tests that include a 375px viewport check

---

## Phase 4 — Deployment & Launch

> TASK-035 and TASK-036 are **parallel**. TASK-037 depends on TASK-035.
> TASK-038 depends on both TASK-036 and TASK-037.

---

#### TASK-035: AWS infrastructure — Terraform
- **Phase:** Phase 4 — Deployment & Launch
- **Depends on:** TASK-031, TASK-032
- **Files:**
  - `infra/terraform/main.tf`
  - `infra/terraform/variables.tf`
  - `infra/terraform/outputs.tf`
  - `infra/terraform/modules/ec2/`
  - `infra/terraform/modules/s3_images/`
  - `infra/terraform/modules/cloudfront_images/`
  - `infra/terraform/modules/route53/`
- **Description:** Write Terraform modules for the single-EC2 production baseline defined in the architecture. Provision: a single EC2 instance (Ubuntu 24.04, appropriately right-sized) with Docker Engine + Docker Compose v2, a gp3 EBS volume for PostgreSQL data and Redis AOF persistence, an S3 bucket (private, no public access) for product images with a CloudFront distribution in front of it, security groups allowing only HTTPS (443) and SSH from a known CIDR, and Route 53 DNS records. Terraform state in S3 + DynamoDB locking. PostgreSQL and Redis run as Docker Compose services on the EC2 host — do NOT provision RDS or ElastiCache; those are future upgrade options documented in the architecture.
- **Acceptance criteria:**
  - [ ] `terraform plan` reports no errors; `terraform apply` provisions all resources
  - [ ] EC2 instance has Docker Engine and Docker Compose v2 installed (via user-data or provisioner)
  - [ ] EBS volume is mounted and used for Postgres data and Redis AOF
  - [ ] S3 bucket blocks all public access; product images are served via CloudFront with presigned URLs
  - [ ] Security group allows only port 443 (HTTPS) and SSH from a restricted CIDR
  - [ ] Terraform state is stored remotely in S3 with DynamoDB locking
  - [ ] No RDS or ElastiCache resources are created (those are out of scope for Phase 4)

> **ADR needed:** Confirmed single-EC2 + Docker Compose baseline; document ECS/RDS as the documented upgrade path once throughput or reliability justify it

---

#### TASK-036: Frontend static export & CDN deployment
- **Phase:** Phase 4 — Deployment & Launch
- **Depends on:** TASK-035
- **Files:**
  - `frontend/next.config.ts` *(verify `output: 'export'` settings)*
  - `infra/terraform/modules/cloudfront_frontend/`
  - `Makefile` *(add `deploy-frontend` target)*
- **Description:** Provision a second S3 bucket + CloudFront distribution for the frontend static export. Add `make deploy-frontend` which runs `next build`, syncs `out/` to S3, and invalidates the CloudFront cache. Configure HTTPS-only with an ACM certificate. Configure CloudFront error pages to serve `404.html` from the export for unknown routes.
- **Acceptance criteria:**
  - [ ] `make deploy-frontend` completes without errors
  - [ ] All pre-rendered routes return `200` via CloudFront
  - [ ] Unknown routes return the custom 404 page
  - [ ] HTTP redirects to HTTPS; HSTS header is present

---

#### TASK-037: Secrets management & production config
- **Phase:** Phase 4 — Deployment & Launch
- **Depends on:** TASK-035
- **Files:**
  - `backend/internal/config/config.go` *(update to read from Secrets Manager in production)*
  - `.env.example` *(final update — all production vars documented)*
- **Description:** Update the backend `Config` loader: when `APP_ENV=production`, read secrets from AWS Secrets Manager instead of environment variables. Document every required secret in `.env.example` with descriptions and example values. No secret should ever be committed to git.
- **Acceptance criteria:**
  - [ ] Backend starts successfully in production mode using Secrets Manager values
  - [ ] No secret value is hardcoded or committed anywhere in the repository
  - [ ] `.env.example` lists every required variable with a description and `CHANGE_ME` placeholder
  - [ ] `make check-env` confirms all required vars are present in the current environment

---

#### TASK-038: Launch checklist & smoke tests
- **Phase:** Phase 4 — Deployment & Launch
- **Depends on:** TASK-036, TASK-037
- **Files:**
  - `verify/launch-checklist.md`
  - `tests/e2e/smoke.spec.ts`
- **Description:** Write a Playwright smoke-test suite that exercises the critical user journey end-to-end against the production URL: landing page → product list → product detail → add to cart → checkout. Create `verify/launch-checklist.md` covering: SSL certificate validity, DNS resolution, Facebook OAuth in production, OTP email delivery, S3 image serving via CloudFront, admin login and seed data verification.
- **Acceptance criteria:**
  - [ ] All smoke tests pass against the production URL
  - [ ] SSL certificate is valid and set to auto-renew (ACM)
  - [ ] Admin can log in with seeded credentials in production
  - [ ] Both seeded coupons (`SUPERHIT`, `FIRSTTIME`) are present and functional in production
  - [ ] Google Maps loads on the Contact Us page in production
  - [ ] `verify/launch-checklist.md` is fully checked off before go-live

---

## Dependency Graph Summary

```
Phase 0 — sequential
  TASK-001 → TASK-002, TASK-003, TASK-004  (002 & 003 in parallel)
  TASK-004 → TASK-040 (worker scaffold)

Phase 1 — TASK-005 & TASK-007 & TASK-008 in parallel after Phase 0
  TASK-004 → TASK-005 → TASK-006
  TASK-003 → TASK-007
  TASK-003 → TASK-008

Phase 2 — all sub-phases parallel once Phase 1 is complete
  2A: TASK-009 → TASK-010
  2B: TASK-005 → TASK-011 → TASK-012 → TASK-013 (sequential)
      TASK-012 + TASK-007 → TASK-014, TASK-015, TASK-016 (parallel)
  2C: TASK-003 → TASK-017, TASK-018, TASK-019 (parallel)
      TASK-003 + TASK-005 → TASK-020
  2D: TASK-006 + TASK-005 → TASK-021 → TASK-022
      TASK-021 → TASK-023 → TASK-039 (Razorpay) → TASK-024
  2E: TASK-006 → TASK-025 → TASK-026, TASK-027, TASK-028, TASK-029, TASK-041 (parallel)
  2F: TASK-005 + TASK-006 → TASK-030

Phase 3 — all parallel after Phase 2 complete
  TASK-031, TASK-032, TASK-033, TASK-034

Phase 4 — sequential pair then parallel
  TASK-035 ∥ TASK-036
  TASK-035 → TASK-037
  TASK-036 + TASK-037 → TASK-038
```

### Phase overlap summary

| Pair | Can overlap? | Notes |
|---|---|---|
| Phase 0 + Phase 1 | Partially | Phase 1 TASK-007/008 can start as soon as TASK-003 is done |
| Phase 1 + Phase 2 | Partially | 2B backend (TASK-011) can start as soon as TASK-005 is done |
| Phase 2 sub-phases | Yes | 2A, 2B, 2C, 2D, 2E, 2F are independent of each other |
| Phase 2 + Phase 3 | Partially | TASK-031/032/033 can begin per-sub-phase as features complete |
| Phase 3 + Phase 4 | Partially | Terraform (TASK-035) can start while hardening is ongoing |

---

**Total tasks: 41**
**ADRs flagged: 5** (TASK-002, TASK-003, TASK-006, TASK-035, TASK-039)
