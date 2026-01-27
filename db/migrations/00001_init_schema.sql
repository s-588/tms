-- +goose Up
-- +goose StatementBegin
create table clients(
    client_id serial primary key,
    name varchar(50) not null,
    email varchar(254) not null,
    email_verified boolean not null default false,
    phone varchar(25) check (phone ~ '^[\\+]?[0-9\\-\\s()]{7,25}$') not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table employees(
    employee_id serial primary key,
    name varchar(50) not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table fuels(
    fuel_id serial primary key,
    name varchar(50) not null,
    supplier varchar(50),
    price money not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table orders(
    order_id serial primary key,
    distance integer not null,
    weight integer not null,
    total_price money not null,
    status varchar(50) not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table prices(
    price_id serial primary key,
    cargo_type varchar(50) not null unique,
    cost money not null,
    weight integer not null,
    distance integer not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table transports(
    transport_id serial primary key,
    employee_id integer references employees(employee_id),
    model varchar(50) not null,
    license_plate varchar(10),
    payload_capacity integer not null,
    fuel_id integer references fuels(fuel_id) not null,
    fuel_consumption integer not null
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table orders_transport(
    order_id integer references orders(order_id),
    transport_id integer references transports(transport_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
create table clients_orders(
    client_id integer references clients(client_id) not null,
    order_id integer references orders(order_id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transports;
-- +goose StatementEnd

-- +goose StatementBegin
drop table clients;
-- +goose StatementEnd

-- +goose StatementBegin
drop table fuels;
-- +goose StatementEnd

-- +goose StatementBegin
drop table orders;
-- +goose StatementEnd

-- +goose StatementBegin
drop table orders_transport;
-- +goose StatementEnd

-- +goose StatementBegin
drop table clients_orders;
-- +goose StatementEnd

-- +goose StatementBegin
drop table employees;
-- +goose StatementEnd

-- +goose StatementBegin
drop table prices;
-- +goose StatementEnd
