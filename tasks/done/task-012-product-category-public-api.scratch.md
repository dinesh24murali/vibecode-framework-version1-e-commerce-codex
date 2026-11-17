# Scratch: task-012 product & category API

## Plan
1. Migration 22: categories.image_url/updated_at, product_images.alt_text/sort_order (spec fields).
2. catalog.sql: Count/List with optional category, full-text, attribute-value filters and sort; batch queries for images, attributes, inventory, active category slugs.
3. models DTOs, ProductRepository/CategoryRepository, ProductService, Product/Category handlers, public routes.
4. Tests: real-DB integration test for filter SQL; contract test with fakes for response shapes, validation, 404/400.

## Uncertainties (resolved)
- Spec params win over plan: `q`, `category[]`, `attr_value_id[]`, `sort_by`, `page`, `limit` (not attributeValueIds/search/per_page). Both `name[]` and `name` forms accepted.
- Abstract Factory for predicates skipped: predicates are static SQL with optional params in sqlc; a factory adds nothing.
- Spec `categories` is an array but DB has one category per product -> array of one.
- stock = on_hand - reserved (min 0). `popular` sort falls back to newest until order data is used (later task).
- Images: URL returned as stored; presigned read URLs arrive with TASK-013.
- Product detail by UUID only (spec PathId); invalid -> 400.
- DTO structs carry //nolint:govet for fieldalignment, since field order mirrors the spec.

## Skipped
- Landing endpoint, frontend pages (TASK-014+), seed data (TASK-030).

## Risks
- Offset pagination on large catalogues; fine at expected size.
