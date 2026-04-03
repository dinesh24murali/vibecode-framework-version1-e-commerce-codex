-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM ('admin', 'customer');

CREATE TABLE users (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    email           CITEXT      NOT NULL UNIQUE,
    password_hash   TEXT,
    full_name       TEXT        NOT NULL DEFAULT '',
    phone_e164      TEXT,
    role            user_role   NOT NULL DEFAULT 'customer',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_role ON users (role);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
DROP EXTENSION IF EXISTS citext;
