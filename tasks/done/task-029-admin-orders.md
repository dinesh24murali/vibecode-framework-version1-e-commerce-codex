# Task: task-029-admin-orders

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin order management: list (status filter, paginated), detail, status update, plus an orders page with a detail modal.

## Acceptance Criteria
- [x] Admin can filter orders by `open`, `fulfilled`, `cancelled` status tabs
- [x] Status can be changed to `fulfilled` or `cancelled` via the modal dropdown
- [x] Order detail modal shows customer name, delivery address, line items, and total
- [x] Table is paginated (20/page)
- [x] Contract tests cover all 3 admin order endpoints (real Postgres)
- [x] `make verify` passes

## Dependencies
- **Tasks:** 025, 010
- **Memory files:** memory/product.md, memory/auth.md
