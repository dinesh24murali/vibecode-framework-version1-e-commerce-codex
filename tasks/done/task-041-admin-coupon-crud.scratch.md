# Scratch: task-041-admin-coupon-crud

## Plan
1. The admin coupon endpoints were missing from the spec: add paths, Coupon/Create/Update schemas and an admin-coupons tag; regenerate the client.
2. sqlc queries (count, full update, delete, redemption count).
3. CouponAdminRepository (separate from the checkout CouponRepository), AdminCouponService, handler, routes.
4. Coupons page, CouponTable, CouponFormModal; dashboard link.
5. Real-DB contract test + mocked e2e.

## Uncertainties / decisions
- Spec wins and had no coupon admin endpoints: added them (like payments in ADR-0006).
- Only the amount matching `type` is stored; the other column is NULL. Switching type requires the new amount.
- Codes are upper-cased, 2-40 chars of A-Z 0-9 _ -, immutable after creation (matches checkout lookup, which upper-cases).
- `coupon_expiry: null` in an update clears the expiry (absent leaves it).
- Delete 409 when coupon_users rows exist.

## What Will Be Skipped
Usage counts in the table, pagination UI (limit=100), description editing in the UI.

## Risks
Expiry date from the UI is stored as end of day UTC.
