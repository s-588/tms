-- +goose Up
-- +goose StatementBegin
create trigger trigger_employees_updated_at
before update
on employees
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_employees_deleted_at
before delete
on employees
for each row
execute function set_employees_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_prices_updated_at on prices;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_prices_deleted_at on prices;
-- +goose StatementEnd
