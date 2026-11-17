# Scratch: task-028-admin-category-crud

## Plan
1. sqlc queries (create/update/delete/paged list/count products); regenerate.
2. CategoryAdminRepository/Service/Handler + routes under AdminMiddleware.
3. Spec: remove the 12-value enum on Category.slug (it would forbid new categories), add 409 to create/update/delete; regenerate client.
4. Categories page + CategoryFormModal (slug auto-fill); dashboard link.
5. Real-DB contract test + mocked e2e.

## Uncertainties / decisions
- Spec enum on Category.slug contradicted "admin can create categories": removed the enum.
- 409 on delete for ANY linked product (products.category_id is NOT NULL FK, so inactive products also block deletion), not only active ones.
- image_url accepted by the API but no UI field yet.

## What Will Be Skipped
Category image upload UI, pagination UI (limit=100).

## Risks
Slug changes alter public filter URLs (?category=slug).
