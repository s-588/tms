-- +goose Up
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

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_insurances_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE insurances
    SET deleted_at = now()
    WHERE insurance_id = OLD.insurance_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_inspections_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE inspections
    SET deleted_at = now()
    WHERE inspection_id = OLD.inspection_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function if exists set_transports_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
drop function if exists set_insurances_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
drop function if exists set_inspections_deleted_at_trigger_func();
-- +goose StatementEnd
