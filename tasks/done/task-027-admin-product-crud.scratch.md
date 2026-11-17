# Scratch: task-027-admin-product-crud

## Plan
1. sqlc admin queries (search list, update incl. slug, soft delete, image replace, inventory upsert, assignment add/remove).
2. AdminProductRepository (tx for create/update), AdminProductService, handler, routes.
3. Spec: add 409 to create/update; regenerate client.
4. putToPresignedUrl in lib/api/client.ts (only place allowed to fetch); uploader, drawer, table, page; dashboard link.
5. Real-DB contract test; mocked e2e.

## Uncertainties / decisions
- products table has one category_id: only the first of category_ids is stored (>=1 required).
- sku is NOT NULL UNIQUE but not in the spec: sku mirrors slug.
- Admin `stock` = inventory on_hand (public stock = on_hand - reserved).
- discount_percent column is INT: fractional discounts are rounded; 0 stored as NULL.
- Image list on update replaces all images (sort_order from request).

## What Will Be Skipped
Attribute-assignment UI in the drawer (API done and tested); pagination UI (limit=100); image reordering/alt text editing.

## Risks
Slug change affects public product URLs; images replaced wholesale.
