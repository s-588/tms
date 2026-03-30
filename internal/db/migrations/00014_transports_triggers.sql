-- +goose Up
-- +goose StatementBegin
create trigger trigger_transports_updated_at
before update
on transports
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_transports_updated_at on transports;
-- +goose StatementEnd