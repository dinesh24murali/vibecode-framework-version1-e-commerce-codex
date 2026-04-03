-- +goose Up
CREATE TABLE product_attribute_assignments (
    product_id         UUID        NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    attribute_value_id UUID        NOT NULL REFERENCES product_attribute_values (id) ON DELETE CASCADE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (product_id, attribute_value_id)
);

-- +goose Down
DROP TABLE IF EXISTS product_attribute_assignments;
