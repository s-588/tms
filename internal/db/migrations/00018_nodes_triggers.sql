-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER trigger_nodes_updated_at
    BEFORE UPDATE ON nodes
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER trigger_nodes_updated_at ON nodes;
-- +goose StatementEnd