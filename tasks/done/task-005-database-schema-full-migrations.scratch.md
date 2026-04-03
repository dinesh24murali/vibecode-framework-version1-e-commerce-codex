# Scratch: TASK-005 — Database schema — full migrations

## Plan

1. **Write 20 migration SQL files** in `backend/migrations/` using goose timestamp format (consistent with existing `20260401000000_init.sql`)
2. **Create sqlc query files** in `backend/db/queries/` — one per domain
3. **Verify `make migrate-up` and `make sqlc-gen`** succeed

---

## Migration Order (dependency-safe)

Files are ordered so that FK targets always exist before FK sources.

| # | Timestamp | File | Notes |
|---|-----------|------|-------|
| 1 | 20260402000001 | `create_users` | standalone; CITEXT extension; role enum or CHECK |
| 2 | 20260402000002 | `create_refresh_tokens` | FK → users |
| 3 | 20260402000003 | `create_consents` | FK → users |
| 4 | 20260402000004 | `create_categories` | standalone |
| 5 | 20260402000005 | `create_product_attributes` | standalone |
| 6 | 20260402000006 | `create_product_attribute_values` | FK → product_attributes |
| 7 | 20260402000007 | `create_products` | FK → categories; generated tsvector; GIN index |
| 8 | 20260402000008 | `create_product_images` | FK → products |
| 9 | 20260402000009 | `create_product_attribute_assignments` | FK → products, product_attribute_values |
| 10 | 20260402000010 | `create_inventory_items` | FK → products; row-lock ready |
| 11 | 20260402000011 | `create_cart` | FK → users, products; unique (user_id, product_id) |
| 12 | 20260402000012 | `create_orders` | FK → users; order_status enum |
| 13 | 20260402000013 | `create_order_items` | FK → orders, products |
| 14 | 20260402000014 | `create_order_address` | FK → orders only (no FK to addresses) |
| 15 | 20260402000015 | `create_payments` | FK → orders; gateway_payload JSONB |
| 16 | 20260402000016 | `create_audit_events` | FK → orders; actor_user_id no FK (allow deleted users) |
| 17 | 20260402000017 | `create_coupons` | standalone; type CHECK flat|percentage |
| 18 | 20260402000018 | `create_coupon_users` | FK → users, coupons; unique (user_id, coupon_id) |
| 19 | 20260402000019 | `create_addresses` | FK → users |
| 20 | 20260402000020 | `create_contact_queries` | standalone |

---

## Schema Details Per Migration

### Migration 1 — users
- Enable `CREATE EXTENSION IF NOT EXISTS citext`
- Create `user_role` enum: `admin`, `customer`
- Columns: `id UUID DEFAULT gen_random_uuid() PK`, `email CITEXT NOT NULL UNIQUE`, `password_hash TEXT`, `full_name TEXT NOT NULL`, `phone_e164 TEXT`, `role user_role NOT NULL DEFAULT 'customer'`, `created_at TIMESTAMPTZ DEFAULT now()`, `updated_at TIMESTAMPTZ DEFAULT now()`, `deleted_at TIMESTAMPTZ`
- Indexes: unique on email (already from UNIQUE constraint), btree on role

### Migration 2 — refresh_tokens
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `token_family_id TEXT NOT NULL`, `token_hash TEXT NOT NULL`, `user_agent TEXT`, `ip_address INET`, `expires_at TIMESTAMPTZ NOT NULL`, `revoked_at TIMESTAMPTZ`, `created_at TIMESTAMPTZ`
- Index: btree on `(user_id, expires_at)`

### Migration 3 — consents
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `consent_type TEXT NOT NULL`, `policy_version TEXT NOT NULL`, `granted_at TIMESTAMPTZ NOT NULL`, `withdrawn_at TIMESTAMPTZ`
- Index: btree on `(user_id, consent_type, granted_at DESC)`

### Migration 4 — categories
- Columns: `id UUID PK`, `slug TEXT NOT NULL UNIQUE`, `name TEXT NOT NULL`, `description TEXT`, `created_at TIMESTAMPTZ`

### Migration 5 — product_attributes
- Columns: `id UUID PK`, `code TEXT NOT NULL UNIQUE`, `name TEXT NOT NULL`, `is_filterable BOOLEAN NOT NULL DEFAULT false`, `sort_order INT NOT NULL DEFAULT 0`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Index: unique btree on code (from constraint)

### Migration 6 — product_attribute_values
- Columns: `id UUID PK`, `attribute_id UUID NOT NULL FK→product_attributes`, `value TEXT NOT NULL`, `slug TEXT NOT NULL`, `sort_order INT NOT NULL DEFAULT 0`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Index: unique btree on `(attribute_id, slug)`

### Migration 7 — products
- Columns: `id UUID PK`, `category_id UUID NOT NULL FK→categories`, `sku TEXT NOT NULL UNIQUE`, `slug TEXT NOT NULL UNIQUE`, `title TEXT NOT NULL`, `author_name TEXT NOT NULL`, `isbn13 TEXT UNIQUE`, `description TEXT`, `price_inr NUMERIC(12,2) NOT NULL`, `discount_percent INT CHECK (discount_percent > 0 AND discount_percent <= 100)`, `image_url TEXT`, `currency_code CHAR(3) NOT NULL DEFAULT 'INR'`, `language_code TEXT`, `page_count INT`, `publication_date DATE`, `is_active BOOLEAN NOT NULL DEFAULT true`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Generated column: `search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', coalesce(title,'') || ' ' || coalesce(author_name,'') || ' ' || coalesce(description,'') || ' ' || coalesce(isbn13,''))) STORED`
- Indexes: unique on sku, slug, isbn13; btree on `(category_id, is_active)`; btree on `(is_active, created_at DESC)`; GIN on `search_vector`

### Migration 8 — product_images
- Columns: `id UUID PK`, `product_id UUID NOT NULL FK→products`, `image_url TEXT NOT NULL`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Index: btree on `(product_id, created_at DESC)`

### Migration 9 — product_attribute_assignments
- Columns: `product_id UUID NOT NULL FK→products`, `attribute_value_id UUID NOT NULL FK→product_attribute_values`, `created_at TIMESTAMPTZ`
- PK: `(product_id, attribute_value_id)` composite
- Index: unique btree on `(product_id, attribute_value_id)` (from PK)

### Migration 10 — inventory_items
- Columns: `product_id UUID PK FK→products`, `on_hand INT NOT NULL DEFAULT 0`, `reserved INT NOT NULL DEFAULT 0`, `reorder_threshold INT NOT NULL DEFAULT 0`, `updated_at TIMESTAMPTZ`
- Note: row-level locking (`SELECT ... FOR UPDATE`) is done in application, not schema

### Migration 11 — cart
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `product_id UUID NOT NULL FK→products`, `quantity INT NOT NULL CHECK (quantity > 0)`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Unique constraint on `(user_id, product_id)`
- Indexes: unique btree on `(user_id, product_id)`; btree on `(user_id, updated_at DESC)`

### Migration 12 — orders
- Create `order_status` enum: `pending_payment`, `paid`, `packed`, `shipped`, `delivered`, `cancelled`, `payment_failed`, `refunded`
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `order_number TEXT NOT NULL UNIQUE`, `status order_status NOT NULL DEFAULT 'pending_payment'`, `subtotal_inr NUMERIC(12,2) NOT NULL`, `shipping_inr NUMERIC(12,2) NOT NULL DEFAULT 0`, `tax_inr NUMERIC(12,2) NOT NULL DEFAULT 0`, `total_inr NUMERIC(12,2) NOT NULL`, `currency_code CHAR(3) NOT NULL DEFAULT 'INR'`, `placed_at TIMESTAMPTZ`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Indexes: btree on `(user_id, created_at DESC)`; unique on order_number

### Migration 13 — order_items
- Columns: `id UUID PK`, `order_id UUID NOT NULL FK→orders`, `product_id UUID FK→products` (nullable — product may be deleted later), `product_title TEXT NOT NULL`, `product_sku TEXT NOT NULL`, `quantity INT NOT NULL`, `unit_price_inr NUMERIC(12,2) NOT NULL`, `line_total_inr NUMERIC(12,2) NOT NULL`

### Migration 14 — order_address
- Columns: `id UUID PK`, `order_id UUID NOT NULL UNIQUE FK→orders`, `recipient_name TEXT NOT NULL`, `line1 TEXT NOT NULL`, `line2 TEXT`, `city TEXT NOT NULL`, `state TEXT NOT NULL`, `postal_code TEXT NOT NULL`, `country_code TEXT NOT NULL`, `phone_e164 TEXT`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Index: unique btree on order_id (from UNIQUE constraint)
- **No FK to addresses table**

### Migration 15 — payments
- Columns: `id UUID PK`, `order_id UUID NOT NULL FK→orders`, `provider TEXT NOT NULL`, `provider_order_id TEXT NOT NULL UNIQUE`, `provider_payment_id TEXT`, `status TEXT NOT NULL`, `amount_inr NUMERIC(12,2) NOT NULL`, `gateway_payload JSONB`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Indexes: unique on provider_order_id; btree on `(order_id, status)`

### Migration 16 — audit_events
- Columns: `id UUID PK`, `order_id UUID FK→orders`, `actor_user_id UUID` (no FK — allow deleted users in audit trail), `event_type TEXT NOT NULL`, `payload JSONB`, `created_at TIMESTAMPTZ`
- Index: btree on `(order_id, created_at DESC)`

### Migration 17 — coupons
- Columns: `id UUID PK`, `name TEXT NOT NULL`, `description TEXT`, `code TEXT NOT NULL UNIQUE`, `is_enable BOOLEAN NOT NULL DEFAULT true`, `type TEXT NOT NULL CHECK (type IN ('flat', 'percentage'))`, `percent INT`, `flat INT`, `coupon_expiry TIMESTAMPTZ`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Index: unique on code

### Migration 18 — coupon_users
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `coupon_id UUID NOT NULL FK→coupons`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ`
- Unique constraint on `(user_id, coupon_id)`
- Indexes: unique btree on `(user_id, coupon_id)`; btree on `(coupon_id, created_at DESC)`

### Migration 19 — addresses
- Columns: `id UUID PK`, `user_id UUID NOT NULL FK→users`, `recipient_name TEXT NOT NULL`, `line1 TEXT NOT NULL`, `line2 TEXT`, `city TEXT NOT NULL`, `state TEXT NOT NULL`, `postal_code TEXT NOT NULL`, `country_code TEXT NOT NULL`, `phone_e164 TEXT`, `is_default BOOLEAN NOT NULL DEFAULT false`, `created_at TIMESTAMPTZ`

### Migration 20 — contact_queries
- Columns: `id UUID PK`, `name TEXT NOT NULL`, `email TEXT NOT NULL`, `phone TEXT`, `subject TEXT`, `message TEXT NOT NULL`, `created_at TIMESTAMPTZ`

---

## sqlc Query Files (`backend/db/queries/`)

| File | Queries |
|------|---------|
| `users.sql` | GetUserByEmail, GetUserByID, CreateUser, UpdateUser, SoftDeleteUser |
| `refresh_tokens.sql` | CreateRefreshToken, GetRefreshTokenByHash, RevokeRefreshToken, RevokeTokenFamily, DeleteExpiredTokens |
| `consents.sql` | CreateConsent, GetConsentsByUser, WithdrawConsent |
| `categories.sql` | ListCategories, GetCategoryBySlug, CreateCategory, UpdateCategory |
| `products.sql` | ListProducts, GetProductBySlug, GetProductByID, SearchProducts, CreateProduct, UpdateProduct |
| `product_images.sql` | ListImagesByProduct, CreateProductImage, DeleteProductImage |
| `product_attributes.sql` | ListAttributes, GetAttributeByCode, CreateAttribute |
| `product_attribute_values.sql` | ListValuesByAttribute, CreateAttributeValue |
| `product_attribute_assignments.sql` | AssignAttributeToProduct, GetAttributesByProduct, RemoveAttributeFromProduct |
| `inventory.sql` | GetInventoryByProduct, UpdateInventory, ReserveStock, ReleaseStock |
| `cart.sql` | GetCartByUser, UpsertCartItem, DeleteCartItem, ClearCart |
| `orders.sql` | CreateOrder, GetOrderByID, GetOrderByNumber, ListOrdersByUser, UpdateOrderStatus |
| `order_items.sql` | CreateOrderItem, ListOrderItems |
| `order_address.sql` | CreateOrderAddress, GetOrderAddress |
| `payments.sql` | CreatePayment, GetPaymentByProviderOrderID, UpdatePaymentStatus |
| `audit_events.sql` | CreateAuditEvent, ListAuditEventsByOrder |
| `coupons.sql` | GetCouponByCode, ListCoupons, CreateCoupon, UpdateCoupon |
| `coupon_users.sql` | CreateCouponUser, GetCouponUserByUserAndCoupon |
| `addresses.sql` | ListAddressesByUser, GetAddressByID, CreateAddress, UpdateAddress, DeleteAddress, SetDefaultAddress |
| `contact_queries.sql` | CreateContactQuery, ListContactQueries |

---

## Files to Change

**Create (migrations):**
- `backend/migrations/20260402000001_create_users.sql`
- `backend/migrations/20260402000002_create_refresh_tokens.sql`
- `backend/migrations/20260402000003_create_consents.sql`
- `backend/migrations/20260402000004_create_categories.sql`
- `backend/migrations/20260402000005_create_product_attributes.sql`
- `backend/migrations/20260402000006_create_product_attribute_values.sql`
- `backend/migrations/20260402000007_create_products.sql`
- `backend/migrations/20260402000008_create_product_images.sql`
- `backend/migrations/20260402000009_create_product_attribute_assignments.sql`
- `backend/migrations/20260402000010_create_inventory_items.sql`
- `backend/migrations/20260402000011_create_cart.sql`
- `backend/migrations/20260402000012_create_orders.sql`
- `backend/migrations/20260402000013_create_order_items.sql`
- `backend/migrations/20260402000014_create_order_address.sql`
- `backend/migrations/20260402000015_create_payments.sql`
- `backend/migrations/20260402000016_create_audit_events.sql`
- `backend/migrations/20260402000017_create_coupons.sql`
- `backend/migrations/20260402000018_create_coupon_users.sql`
- `backend/migrations/20260402000019_create_addresses.sql`
- `backend/migrations/20260402000020_create_contact_queries.sql`

**Create (sqlc queries):**
- `backend/db/queries/users.sql`
- `backend/db/queries/refresh_tokens.sql`
- `backend/db/queries/consents.sql`
- `backend/db/queries/categories.sql`
- `backend/db/queries/products.sql`
- `backend/db/queries/product_images.sql`
- `backend/db/queries/product_attributes.sql`
- `backend/db/queries/product_attribute_values.sql`
- `backend/db/queries/product_attribute_assignments.sql`
- `backend/db/queries/inventory.sql`
- `backend/db/queries/cart.sql`
- `backend/db/queries/orders.sql`
- `backend/db/queries/order_items.sql`
- `backend/db/queries/order_address.sql`
- `backend/db/queries/payments.sql`
- `backend/db/queries/audit_events.sql`
- `backend/db/queries/coupons.sql`
- `backend/db/queries/coupon_users.sql`
- `backend/db/queries/addresses.sql`
- `backend/db/queries/contact_queries.sql`

---

## Uncertainties

1. **Goose file naming:** The existing migration uses `20260401000000_init.sql` (timestamp format). The task file describes `0001_create_users.sql` (sequential format). Using timestamp format for consistency with what was scaffolded in TASK-004. Goose supports both, but cannot mix within a repo — sticking with timestamp.

2. **`audit_events.actor_user_id` FK:** Keeping it as a bare UUID with no FK so audit records survive user deletion. Noted explicitly.

3. **`coupons.type` vs architecture schema:** Architecture schema shows `percent` and `flat` as separate INT columns with no `type` discriminator column, but the task acceptance criteria requires `type TEXT CHECK (type IN ('flat', 'percentage'))`. Adding the `type` column as required by the task — this is additive and safe.

4. **`make sqlc-gen` runs against a live DB:** sqlc generation only needs the migration SQL files (schema source: `migrations/`), not a running DB. Should succeed offline.

5. **`product_attribute_assignments` table name:** Architecture uses `PRODUCT_ATTRIBUTE_ASSIGNMENTS`, task file uses same. Mapping to snake_case: `product_attribute_assignments`.

---

## What Will Be Skipped

- Seed data migrations — kept separate per architecture migration policy (§3)
- Enum migrations for `order_status` and `user_role`: created inline within the relevant migration file (migrations 1 and 12 respectively), not as separate migration files — this is standard practice
- No Go code changes in this task — sqlc generates `backend/db/sqlc/` automatically; that directory will be updated by `make sqlc-gen` but is not hand-edited

---

## Risks

1. **`GENERATED ALWAYS AS ... STORED` on `search_vector`:** Supported in PostgreSQL 12+; our target is PostgreSQL 18 so this is safe. The generated column syntax is correct for PG.
2. **sqlc compatibility with generated columns:** sqlc v1.30.x supports PostgreSQL generated columns — the `search_vector` column must be excluded from INSERT/UPDATE queries. Will mark queries accordingly with `-- name: ... :many` annotations.
3. **Down migrations must drop in reverse dependency order:** each `-- +goose Down` block must drop objects in correct order. Will use `DROP TABLE IF EXISTS CASCADE` where needed.
4. **CITEXT extension:** Must be created before the `users` table. Placing `CREATE EXTENSION IF NOT EXISTS citext` at the top of migration 1.
