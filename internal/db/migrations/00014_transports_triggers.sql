-- +goose Up
-- +goose StatementBegin
create trigger trigger_transports_updated_at
before update
on transports
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_transports_deleted_at
before delete
on transports
for each row
execute function set_transports_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_transports_updated_at on transports;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_transports_deleted_at on transports;
-- +goose StatementEnd
