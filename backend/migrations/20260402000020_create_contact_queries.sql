-- +goose Up
CREATE TABLE contact_queries (
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    phone      TEXT,
    subject    TEXT,
    message    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS contact_queries;
