-- +goose Up
CREATE TABLE products (
    id               UUID           NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    category_id      UUID           NOT NULL REFERENCES categories (id),
    sku              TEXT           NOT NULL UNIQUE,
    slug             TEXT           NOT NULL UNIQUE,
    title            TEXT           NOT NULL,
    author_name      TEXT           NOT NULL DEFAULT '',
    isbn13           TEXT           UNIQUE,
    description      TEXT,
    price_inr        NUMERIC(12, 2) NOT NULL,
    discount_percent INT            CHECK (discount_percent > 0 AND discount_percent <= 100),
    image_url        TEXT,
    currency_code    CHAR(3)        NOT NULL DEFAULT 'INR',
    language_code    TEXT,
    page_count       INT,
    publication_date DATE,
    is_active        BOOLEAN        NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    search_vector    TSVECTOR       GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            coalesce(title, '') || ' ' ||
            coalesce(author_name, '') || ' ' ||
            coalesce(description, '') || ' ' ||
            coalesce(isbn13, '')
        )
    ) STORED
);

CREATE INDEX idx_products_category_active   ON products (category_id, is_active);
CREATE INDEX idx_products_active_created    ON products (is_active, created_at DESC);
CREATE INDEX idx_products_search_vector_gin ON products USING GIN (search_vector);

-- +goose Down
DROP TABLE IF EXISTS products;
