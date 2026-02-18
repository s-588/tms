-- +goose Up
-- +goose StatementBegin
create trigger trigger_nodes_updated_at
before update
on nodes
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_nodes_deleted_at
before delete
on nodes
for each row
execute function set_nodes_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_nodes_updated_at on prices;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_nodes_deleted_at on prices;
-- +goose StatementEnd
