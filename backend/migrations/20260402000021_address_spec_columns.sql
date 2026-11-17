-- +goose Up
-- Align addresses with the API spec Address schema (landmark, company, GST, updated_at)
-- and let spec-driven writes omit recipient_name / country_code.
ALTER TABLE addresses
    ADD COLUMN landmark     TEXT,
    ADD COLUMN company_name TEXT,
    ADD COLUMN gst_number   TEXT,
    ADD COLUMN updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    ALTER COLUMN recipient_name SET DEFAULT '',
    ALTER COLUMN country_code   SET DEFAULT 'IN';

-- +goose Down
ALTER TABLE addresses
    DROP COLUMN landmark,
    DROP COLUMN company_name,
    DROP COLUMN gst_number,
    DROP COLUMN updated_at,
    ALTER COLUMN recipient_name DROP DEFAULT,
    ALTER COLUMN country_code   DROP DEFAULT;
