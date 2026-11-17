-- +goose Up
-- Align categories and product images with the API spec (Category.image_url/updated_at,
-- ProductImage.alt_text/sort_order).
ALTER TABLE categories
    ADD COLUMN image_url  TEXT,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

ALTER TABLE product_images
    ADD COLUMN alt_text   TEXT,
    ADD COLUMN sort_order INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE product_images
    DROP COLUMN alt_text,
    DROP COLUMN sort_order;

ALTER TABLE categories
    DROP COLUMN image_url,
    DROP COLUMN updated_at;
