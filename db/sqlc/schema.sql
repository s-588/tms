create table clients(
    client_ID serial primary key,
    name varchar(50) not null,
    email varchar(254) not null,
    phone varchar(25) check (phone ~ '^[\\+]?[0-9\\-\\s()]{7,25}$') not null
);
create table employees(
    employee_id serial primary key,
    name varchar(50) not null
);
create table fuel_types(
    fuel_id serial primary key,
    name varchar(50) not null,
    supplier varchar(50),
    price money not null
);
create table orders(
    order_id serial primary key,
    distance integer not null,
    weight integer not null,
    total_price money not null
);
create table prices(
    price_id serial primary key,
    cargo_type varchar(50) not null unique,
    cost money not null,
    weight integer not null,
    distance integer not null
);
create table transports(
    transport_id serial primary key,
    employee_id integer references employees(employee_id),
    model varchar(50) not null,
    license_plate varchar(10),
    payload_capacity integer not null,
    fuel_id integer references fuel_types(fuel_id) not null,
    fuel_consumption integer not null
);
create table orders_transport(
    order_id integer references orders(order_id),
    transport_id integer references transports(transport_id)
);
create table clients_orders(
    client_id integer references clients(client_id) not null,
    order_id integer references orders(order_id) not null
);
