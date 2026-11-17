# Task: task-027-admin-product-crud

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin product CRUD (8 endpoints), soft delete, attribute assignment API, and a modal/drawer products page with S3 image upload.

## Acceptance Criteria
- [x] Admin can create, edit, soft-delete products from the list page via modal
- [x] "Edit" row action opens the modal pre-populated via `GET /api/v1/admin/products/{id}`
- [x] Image upload uses the presigned PUT URL; the public URL is saved with the product
- [x] Attribute assignment (`POST /api/v1/admin/products/{id}/attributes`) links attribute values to the product (API + contract test)
- [x] Attribute removal (`DELETE .../attributes/{attrId}`) unlinks a single value (API + contract test)
- [x] Soft-deleted products are excluded from public endpoints
- [x] No `[id]` dynamic route segment exists under `app/(admin)/admin/products/`
- [x] Contract tests cover all 8 admin product endpoints (real Postgres)
- [x] `make verify` passes

## Dependencies
- **Tasks:** 025, 012, 013
- **Memory files:** memory/product.md, memory/auth.md
