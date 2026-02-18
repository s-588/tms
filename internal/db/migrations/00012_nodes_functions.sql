-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_nodes_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE nodes
    SET deleted_at = now()
    WHERE node_id = OLD.node_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function set_nodes_deleted_at_trigger_func();
-- +goose StatementEnd
