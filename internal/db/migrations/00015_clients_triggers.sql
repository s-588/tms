-- +goose Up
-- +goose StatementBegin
create trigger trigger_clients_updated_at
before update
on clients
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_clients_deleted_at
before delete
on clients
for each row
execute function set_clients_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_clients_updated_at on clients;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_clients_deleted_at on clients;
-- +goose StatementEnd
