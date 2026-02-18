-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_orders_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET deleted_at = now()
    WHERE order_id = OLD.order_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_orders_deleted_at_trigger_func();
-- +goose StatementEnd
