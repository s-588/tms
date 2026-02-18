-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_prices_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE prices
    SET deleted_at = now()
    WHERE price_id = OLD.price_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_prices_deleted_at_trigger_func();
-- +goose StatementEnd
