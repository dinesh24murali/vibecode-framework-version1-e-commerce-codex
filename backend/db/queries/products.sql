-- name: ListProducts :many
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE category_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetProductBySlug :one
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE slug = $1 AND is_active = true;

-- name: GetProductByID :one
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE id = $1;

-- name: SearchProducts :many
SELECT id, category_id, sku, slug, title, author_name, isbn13, description,
       price_inr, discount_percent, image_url, currency_code, language_code,
       page_count, publication_date, is_active, created_at, updated_at
FROM products
WHERE is_active = true
  AND search_vector @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC
LIMIT $2 OFFSET $3;

-- name: CreateProduct :one
INSERT INTO products (
    category_id, sku, slug, title, author_name, isbn13, description,
    price_inr, discount_percent, image_url, currency_code, language_code,
    page_count, publication_date, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12,
    $13, $14, $15
)
RETURNING id, category_id, sku, slug, title, author_name, isbn13, description,
          price_inr, discount_percent, image_url, currency_code, language_code,
          page_count, publication_date, is_active, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products
SET category_id      = $2,
    title            = $3,
    author_name      = $4,
    description      = $5,
    price_inr        = $6,
    discount_percent = $7,
    image_url        = $8,
    is_active        = $9,
    updated_at       = now()
WHERE id = $1
RETURNING id, category_id, sku, slug, title, author_name, isbn13, description,
          price_inr, discount_percent, image_url, currency_code, language_code,
          page_count, publication_date, is_active, created_at, updated_at;
