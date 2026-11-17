-- Admin product management (includes inactive products).

-- name: AdminCountProducts :one
SELECT count(*) FROM products p
WHERE (sqlc.narg(query)::text IS NULL
       OR p.title ILIKE '%' || sqlc.narg(query)::text || '%'
       OR p.slug  ILIKE '%' || sqlc.narg(query)::text || '%'
       OR p.sku   ILIKE '%' || sqlc.narg(query)::text || '%');

-- name: AdminListProducts :many
SELECT p.id, p.category_id, p.sku, p.slug, p.title, p.author_name, p.isbn13, p.description,
       p.price_inr, p.discount_percent, p.image_url, p.currency_code, p.language_code,
       p.page_count, p.publication_date, p.is_active, p.created_at, p.updated_at
FROM products p
WHERE (sqlc.narg(query)::text IS NULL
       OR p.title ILIKE '%' || sqlc.narg(query)::text || '%'
       OR p.slug  ILIKE '%' || sqlc.narg(query)::text || '%'
       OR p.sku   ILIKE '%' || sqlc.narg(query)::text || '%')
ORDER BY p.created_at DESC, p.id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: AdminUpdateProduct :one
UPDATE products
SET category_id      = $2,
    slug             = $3,
    title            = $4,
    description      = $5,
    price_inr        = $6,
    discount_percent = $7,
    is_active        = $8,
    updated_at       = now()
WHERE id = $1
RETURNING id, category_id, sku, slug, title, author_name, isbn13, description,
          price_inr, discount_percent, image_url, currency_code, language_code,
          page_count, publication_date, is_active, created_at, updated_at;

-- name: SetProductActive :execrows
UPDATE products SET is_active = $2, updated_at = now() WHERE id = $1;

-- name: DeleteImagesByProductID :exec
DELETE FROM product_images WHERE product_id = $1;

-- name: InsertProductImage :exec
INSERT INTO product_images (product_id, image_url, alt_text, sort_order)
VALUES ($1, $2, $3, $4);

-- name: UpsertInventoryOnHand :exec
INSERT INTO inventory_items (product_id, on_hand)
VALUES ($1, $2)
ON CONFLICT (product_id) DO UPDATE SET on_hand = EXCLUDED.on_hand, updated_at = now();

-- name: RemoveProductAttributeValue :execrows
DELETE FROM product_attribute_assignments
WHERE product_id = $1 AND attribute_value_id = $2;

-- name: AssignProductAttributeValue :exec
INSERT INTO product_attribute_assignments (product_id, attribute_value_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;
