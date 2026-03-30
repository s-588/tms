-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN distance TYPE double precision
    USING distance::double precision;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN distance TYPE integer
    USING distance::integer;  -- Truncates any fractional part
-- +goose StatementEnd