# Task: task-028-admin-category-crud

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin category CRUD (5 endpoints) and a modal-based categories page.

## Acceptance Criteria
- [x] Admin can create, rename, and delete categories
- [x] Deleting a category with linked products returns `409` with a clear message
- [x] Category changes are immediately reflected in the list without a page reload
- [x] No `[id]` dynamic route exists under `app/(admin)/admin/categories/`
- [x] Contract tests cover all 5 category endpoints (real Postgres)
- [x] `make verify` passes

## Dependencies
- **Tasks:** 025, 012
- **Memory files:** memory/product.md
