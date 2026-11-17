# Task: task-026-admin-attributes-crud

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin CRUD for attribute definitions and values (9 endpoints) plus a modal-based admin page.

## Acceptance Criteria
- [x] Admin can create, rename, and delete attribute definitions
- [x] Admin can add, rename, and delete values per definition
- [x] Deleting an attribute value that is assigned to products returns `409 Conflict`
- [x] All CRUD operations work via modals; no `[id]` dynamic route exists
- [x] Contract tests cover all 9 attribute endpoints (backend/tests/contract/admin_attributes_test.go, real Postgres)
- [x] `make verify` passes

## Dependencies
- **Tasks:** 025, 011
- **Memory files:** memory/auth.md, memory/product.md

## Notes
Spec had the value delete cascading assignments; the plan requires 409. Spec updated (reusable `Conflict` response, 409 on create/update/delete) and client regenerated.
