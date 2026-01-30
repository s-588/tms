-- +goose Up
-- +goose StatementBegin
create trigger trigger_clients_updated_at
after update
on clients
for each row
execute function set_updated_at();

create trigger trigger_clients_deleted_at
instead of delete
on clients
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_employees_updated_at
after update
on employees
for each row
execute function set_updated_at();

create trigger trigger_employees_deleted_at
instead of delete
on employees
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_fuels_updated_at
after update
on fuels
for each row
execute function set_updated_at();

create trigger trigger_fuels_deleted_at
instead of delete
on fuels
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_orders_updated_at
after update
on orders
for each row
execute function set_updated_at();

create trigger trigger_orders_deleted_at
instead of delete
on orders
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_prices_updated_at
after update
on prices
for each row
execute function set_updated_at();

create trigger trigger_prices_deleted_at
instead of delete
on prices
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_transports_updated_at
after update
on transports
for each row
execute function set_updated_at();

create trigger trigger_transports_deleted_at
instead of delete
on transports
for each row
execute function set_deleted_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_clients_updated_at;
drop trigger trigger_clients_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_employees_updated_at;
drop trigger trigger_employees_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_fuels_updated_at;
drop trigger trigger_fuels_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_orders_updated_at;
drop trigger trigger_orders_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_prices_updated_at;
drop trigger trigger_prices_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_transports_updated_at;
drop trigger trigger_transports_deleted_at;
-- +goose StatementEnd
