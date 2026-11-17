-- +goose Up
CREATE TABLE carousel_slides (
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    title      TEXT        NOT NULL,
    subtitle   TEXT,
    image_url  TEXT        NOT NULL,
    cta_label  TEXT,
    cta_url    TEXT,
    sort_order INT         NOT NULL DEFAULT 0,
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_carousel_slides_active_order ON carousel_slides (is_active, sort_order);

-- +goose Down
DROP TABLE IF EXISTS carousel_slides;
