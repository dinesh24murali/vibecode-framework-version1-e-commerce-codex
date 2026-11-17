# Task: task-041-admin-coupon-crud

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin coupon CRUD (5 endpoints) and a modal-based coupons page with enable/disable toggle.

## Acceptance Criteria
- [x] Admin can create, edit, enable/disable, and delete coupons
- [x] Coupon `type` field correctly gates whether `flat` or `percent` amount is used
- [x] Deleting a redeemed coupon returns `409` with a clear message
- [x] No `[id]` dynamic route exists under `app/(admin)/admin/coupons/`
- [x] Contract tests cover all 5 admin coupon endpoints (real Postgres)
- [x] `make verify` passes

## Dependencies
- **Tasks:** 025, 023
- **Memory files:** memory/product.md
