-- +goose Up
-- +goose StatementBegin
create trigger trigger_clients_updated_at
before update
on clients
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_clients_deleted_at
before delete
on clients
for each row
execute function set_clients_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_employees_updated_at
after update
on employees
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_employees_deleted_at
before delete
on employees
for each row
execute function set_employees_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_fuels_updated_at
after update
on fuels
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_fuels_deleted_at
before delete
on fuels
for each row
execute function set_fuels_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_orders_updated_at
after update
on orders
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_orders_deleted_at
before delete
on orders
for each row
execute function set_orders_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_prices_updated_at
after update
on prices
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_prices_deleted_at
before delete
on prices
for each row
execute function set_prices_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_transports_updated_at
after update
on transports
for each row
execute function set_updated_at_trigger_func();

create trigger trigger_transports_deleted_at
before delete
on transports
for each row
execute function set_transports_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_clients_updated_at on clients;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_clients_deleted_at on clients;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_employees_updated_at on employees;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_employees_deleted_at on employees;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_fuels_updated_at on fuels;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_fuels_deleted_at on fuels;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_orders_updated_at on orders;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_orders_deleted_at on orders;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_prices_updated_at on prices;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_prices_deleted_at on prices;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_transports_updated_at on transports;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger trigger_transports_deleted_at on transports;
-- +goose StatementEnd
