-- +goose Up
-- +goose StatementBegin
create function set_updated_at_trigger_func()
returns trigger as $$
begin
    NEW.updated_at = now();
    return NEW;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- We need many delete functions because our tables dont have unified 'id' column.

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_clients_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE clients
    SET deleted_at = now()
    WHERE client_id = OLD.client_id;
    RETURN null;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_employees_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE employees
    SET deleted_at = now()
    WHERE employee_id = OLD.employee_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_fuels_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE fuels
    SET deleted_at = now()
    WHERE fuel_id = OLD.fuel_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

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

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_transports_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE transports
    SET deleted_at = now()
    WHERE transport_id = OLD.transport_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function if exists set_updated_at_trigger_func;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_clients_deleted_at_trigger_func();
DROP FUNCTION IF EXISTS set_employees_deleted_at_trigger_func();
DROP FUNCTION IF EXISTS set_fuels_deleted_at_trigger_func();
DROP FUNCTION IF EXISTS set_orders_deleted_at_trigger_func();
DROP FUNCTION IF EXISTS set_prices_deleted_at_trigger_func();
DROP FUNCTION IF EXISTS set_transports_deleted_at_trigger_func();
-- +goose StatementEnd
