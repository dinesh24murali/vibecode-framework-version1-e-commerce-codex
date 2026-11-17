-- Public catalogue reads. Filters (all optional):
--   category_slugs  OR-combined category slugs
--   query           full-text search on title/author/description/isbn
--   value_ids       attribute value ids: OR within one attribute dimension, AND across dimensions
--   sort            price_asc | price_desc | newest | popular (popular falls back to newest until sales data exists)

-- name: CountCatalogProducts :one
SELECT count(*) FROM products p
WHERE p.is_active = true
  AND (cardinality(sqlc.arg(category_slugs)::text[]) = 0
       OR p.category_id IN (SELECT c.id FROM categories c WHERE c.slug = ANY(sqlc.arg(category_slugs)::text[])))
  AND (sqlc.narg(query)::text IS NULL OR p.search_vector @@ plainto_tsquery('english', sqlc.narg(query)::text))
  AND NOT EXISTS (
      SELECT 1
      FROM (SELECT DISTINCT v.attribute_id FROM product_attribute_values v WHERE v.id = ANY(sqlc.arg(value_ids)::uuid[])) dims
      WHERE NOT EXISTS (
          SELECT 1
          FROM product_attribute_assignments paa
          JOIN product_attribute_values v2 ON v2.id = paa.attribute_value_id
          WHERE paa.product_id = p.id
            AND v2.attribute_id = dims.attribute_id
            AND v2.id = ANY(sqlc.arg(value_ids)::uuid[])
      )
  );

-- name: ListCatalogProducts :many
SELECT p.id, p.category_id, p.sku, p.slug, p.title, p.author_name, p.isbn13, p.description,
       p.price_inr, p.discount_percent, p.image_url, p.currency_code, p.language_code,
       p.page_count, p.publication_date, p.is_active, p.created_at, p.updated_at
FROM products p
WHERE p.is_active = true
  AND (cardinality(sqlc.arg(category_slugs)::text[]) = 0
       OR p.category_id IN (SELECT c.id FROM categories c WHERE c.slug = ANY(sqlc.arg(category_slugs)::text[])))
  AND (sqlc.narg(query)::text IS NULL OR p.search_vector @@ plainto_tsquery('english', sqlc.narg(query)::text))
  AND NOT EXISTS (
      SELECT 1
      FROM (SELECT DISTINCT v.attribute_id FROM product_attribute_values v WHERE v.id = ANY(sqlc.arg(value_ids)::uuid[])) dims
      WHERE NOT EXISTS (
          SELECT 1
          FROM product_attribute_assignments paa
          JOIN product_attribute_values v2 ON v2.id = paa.attribute_value_id
          WHERE paa.product_id = p.id
            AND v2.attribute_id = dims.attribute_id
            AND v2.id = ANY(sqlc.arg(value_ids)::uuid[])
      )
  )
ORDER BY
  CASE WHEN sqlc.arg(sort)::text = 'price_asc'  THEN p.price_inr * (100 - COALESCE(p.discount_percent, 0)) END ASC,
  CASE WHEN sqlc.arg(sort)::text = 'price_desc' THEN p.price_inr * (100 - COALESCE(p.discount_percent, 0)) END DESC,
  p.created_at DESC, p.id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: GetActiveProductByID :one
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE id = $1 AND is_active = true;

-- name: ListImagesByProducts :many
SELECT * FROM product_images
WHERE product_id = ANY($1::uuid[])
ORDER BY product_id, sort_order, created_at;

-- name: ListAttributesByProducts :many
SELECT paa.product_id,
       pav.id        AS attribute_value_id,
       pa.id         AS attribute_id,
       pa.code       AS code,
       pa.name       AS attribute_name,
       pav.value     AS value,
       pav.slug      AS value_slug,
       pa.is_filterable AS is_filterable
FROM product_attribute_assignments paa
JOIN product_attribute_values pav ON pav.id = paa.attribute_value_id
JOIN product_attributes pa ON pa.id = pav.attribute_id
WHERE paa.product_id = ANY($1::uuid[])
ORDER BY paa.product_id, pa.sort_order, pav.sort_order;

-- name: ListInventoryByProducts :many
SELECT * FROM inventory_items
WHERE product_id = ANY($1::uuid[]);

-- name: ListActiveCategorySlugs :many
SELECT DISTINCT c.slug FROM categories c
JOIN products p ON p.category_id = c.id AND p.is_active = true
ORDER BY c.slug;
