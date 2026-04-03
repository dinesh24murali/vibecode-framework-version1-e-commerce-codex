# Task: TASK-005 — Database schema — full migrations

**Status:** active
**Created:** 2026-03-31
**ADR refs:** none (references ADR from TASK-004 for sqlc/goose)

---

## Goal

Write all 20 SQL migration files for goose that establish the complete production schema, then run sqlc to generate the corresponding Go query code.

## Background

This task translates the full data model from `docs/02_outputs/03_tech_architecture.md` §3 into goose migration files. Every table, constraint, index, and enum described in the architecture must be captured here. TASK-006 (auth) and all Phase 2 feature tasks depend on this schema being in place.

Key design constraints from the architecture:
- `users.email` is `CITEXT` (case-insensitive unique)
- `products` has a generated `tsvector` column with a GIN index for full-text search
- `order_address` is an immutable snapshot — no FK to `addresses`
- `orders.status` is a Postgres enum
- `inventory_items` mutations require row-level locking
- `consents` is part of the DPDPA compliance baseline

## Acceptance Criteria

- [ ] All 20 migrations apply cleanly with `make migrate-up`
- [ ] Each migration file has valid `-- +goose Up` and `-- +goose Down` blocks
- [ ] All FK constraints are enforced and verified with `psql \d`
- [ ] `products` has a GIN index on the generated `tsvector` search column (combining `title`, `author_name`, `description`, `isbn13`)
- [ ] `product_attribute_assignments` join table correctly links products to attribute values
- [ ] `inventory_items` has `on_hand`, `reserved`, `reorder_threshold` columns
- [ ] `order_address` has no FK referencing the `addresses` table
- [ ] `orders.status` uses a Postgres enum; invalid values are rejected at the DB level
- [ ] `refresh_tokens` has an index on `(user_id, expires_at)`
- [ ] `payments` table has a `gateway_payload JSONB` column
- [ ] `coupons` has a `type` column constrained to `flat | percentage`
- [ ] `coupon_users` has a unique constraint on `(user_id, coupon_id)`
- [ ] `make sqlc-gen` succeeds after all migrations are applied
- [ ] Corresponding `.sql` query files exist in `backend/db/queries/` for each domain

## Dependencies

- **Tasks:** TASK-004 (goose and sqlc tooling must be configured)
- **Memory files:** `memory/infra.md` (pool and Redis config context)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §3 (Data Architecture — full schema + indexing strategy), `docs/02_outputs/04_api_spec.yaml` (schema shapes used by API responses)

## Files Expected to Change

- `backend/migrations/0001_create_users.sql` *(create)*
- `backend/migrations/0002_create_refresh_tokens.sql` *(create)*
- `backend/migrations/0003_create_consents.sql` *(create)*
- `backend/migrations/0004_create_categories.sql` *(create)*
- `backend/migrations/0005_create_product_attributes.sql` *(create)*
- `backend/migrations/0006_create_product_attribute_values.sql` *(create)*
- `backend/migrations/0007_create_products.sql` *(create)*
- `backend/migrations/0008_create_product_images.sql` *(create)*
- `backend/migrations/0009_create_product_attribute_assignments.sql` *(create)*
- `backend/migrations/0010_create_inventory_items.sql` *(create)*
- `backend/migrations/0011_create_cart.sql` *(create)*
- `backend/migrations/0012_create_orders.sql` *(create)*
- `backend/migrations/0013_create_order_items.sql` *(create)*
- `backend/migrations/0014_create_order_address.sql` *(create)*
- `backend/migrations/0015_create_payments.sql` *(create)*
- `backend/migrations/0016_create_audit_events.sql` *(create)*
- `backend/migrations/0017_create_coupons.sql` *(create)*
- `backend/migrations/0018_create_coupon_users.sql` *(create)*
- `backend/migrations/0019_create_addresses.sql` *(create)*
- `backend/migrations/0020_create_contact_queries.sql` *(create)*
- `backend/db/queries/*.sql` *(create — one file per domain)*

## Notes

- `users.role` must be a CHECK constraint or Postgres enum restricted to `admin | customer`
- `products.discount_percent` must allow NULL (no discount) and have a CHECK `> 0 AND <= 100` when not null
- `orders.status` enum values: `pending_payment`, `paid`, `packed`, `shipped`, `delivered`, `cancelled`, `payment_failed`, `refunded`
- `cart` must have a unique constraint on `(user_id, product_id)` — one row per user-product pair
- Full-text search: generated `tsvector` column using `to_tsvector('english', coalesce(title,'') || ' ' || coalesce(author_name,'') || ...)` with a GIN index
- See `docs/02_outputs/03_tech_architecture.md` §3 Indexing Strategy for the full list of required indexes — implement all of them
- Seed data migrations are kept separate from structural schema migrations (per architecture migration policy)
- The `-- +goose Down` block must cleanly reverse the `-- +goose Up` block; use `DROP TABLE IF EXISTS` with dependency ordering

---

## Scratch File

AI: before implementing, create `tasks/active/task-005-database-schema-full-migrations.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
