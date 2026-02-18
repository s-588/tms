-- +goose Up
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

-- +goose Down
-- +goose StatementBegin
drop function if exists set_clients_deleted_at_trigger_func();
-- +goose StatementEnd
