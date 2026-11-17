# Scratch: task-026-admin-attributes-crud

## Plan
1. sqlc queries (paged list, count, update, delete, assignment count); regenerate.
2. AttributeAdminRepository/Service/Handler; routes under AdminMiddleware group.
3. Spec: add Conflict response + 409s; run orval.
4. Admin page + AttributeTable + two modals; dashboard link.
5. Real-DB contract test, e2e spec with mocked API.

## Files to Create or Modify
See git diff.

## Uncertainties
- Spec (cascade) vs plan (409) on value delete: chose 409 so product data is never silently unlinked; spec amended.
- Definition delete still cascades (spec, plan silent).
- Code is immutable after creation (UpdateAttributeDefinitionRequest has no code).

## What Will Be Skipped
Pagination UI (page fetches limit=100); sort_order editing UI; sqlc version header churn (v1.30 -> v1.31.1) is accepted.

## Risks
Admin session is memory-only: e2e navigates client-side.
