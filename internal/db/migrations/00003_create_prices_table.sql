-- +goose Up
-- +goose StatementBegin
CREATE TABLE prices (
    price_id    SERIAL PRIMARY KEY,
    cargo_type  VARCHAR(50) NOT NULL,
    weight      INTEGER NOT NULL,
    distance    INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    UNIQUE (cargo_type, weight, distance)
);

CREATE INDEX idx_prices_deleted_at
    ON prices (deleted_at)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prices;
-- +goose StatementEnd
