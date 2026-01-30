-- +goose Up
-- +goose StatementBegin
create index idx_clients_client_id on clients(client_id);
-- +goose StatementEnd

-- +goose StatementBegin
create index idx_employees_employee_id on employees(employee_id);
-- +goose StatementEnd

-- +goose StatementBegin
create index idx_fuels_fuel_id on fuels(fuel_id);
-- +goose StatementEnd

-- +goose StatementBegin
create index idx_orders_order_id on orders(order_id);
-- +goose StatementEnd

-- +goose StatementBegin
create index idx_prices_price_id on prices(price_id);
-- +goose StatementEnd

-- +goose StatementBegin
create index idx_transports_transport_id on transports(transport_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_clients_client_id;
-- +goose StatementEnd

-- +goose StatementBegin
drop index idx_employees_employee_id;
-- +goose StatementEnd

-- +goose StatementBegin
drop index idx_fuels_fuel_id;
-- +goose StatementEnd

-- +goose StatementBegin
drop index idx_orders_order_id;
-- +goose StatementEnd

-- +goose StatementBegin
drop index idx_prices_price_id;
-- +goose StatementEnd

-- +goose StatementBegin
drop index idx_transports_transport_id;
-- +goose StatementEnd
