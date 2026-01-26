# Database schema

## Clients table
``` create table clients(
    client_ID serial primary key,
    name varchar(50) not null,
    email varchar(254) not null,
    phone varchar(25) check (phone regexp '^[+]?[0-9\\-\\s()]{7,25}$') not null
);```
* name - clients full name
* email - email address for notifications.
* phone - phone number for extra notifications.
---

## Transports table
```
create table transports(
    transport_id serial primary key,
    employee_id integer references employees(employee_id),
    model varchar(50) not null,
    license_plate varchar(10),
    payload_capacity integer not null,
    fuel_id integer references fuel_types(fuel_id) not null,
    fuel_consumption integer not null
);```
* employee_id - who is currently responsible for this vehicle.
* model - vehicle model name.
* license_plate - vehicle registation plate number.
* payload_capacity - maximum weight vehicle can deliver.
* fuel_id - fuel type.
* fuel_consumption - fuel consumption per 100 km.
---

## Fuel types table

``` create table fuel_types(
    fuel_id serial primary key,
    name varchar(50) not null,
    supplier varchar(50),
    price money not null
);```
* name - fuel type name.
* supplier - supplier of the fuel.
* price - fuel price per 1 liter.
---

## Orders table
``` create table orders(
    order_id serial primary key,
    distance integer not null,
    weight integer not null,
    total_price money not null
); ```
* distance - distance between cargo loading and end delivery points.
* total_price - total price client will pay.
Total price calculated as: fuel_full_price + (weight * distance * prices_total_coefficient).
Where:
    - fuel_full_price = fuel(price) * orders(distance) * transport(fuel_consumption)
    - prices_total_coefficient = prices(cost) + prices(weight) + prices(distance)
---

## Orders transport junction table
One order can be transported by multiple vehicles.
```create table orders_transport(
    order_id integer references orders(order_id),
    transport_id integer references transports(transport_id)
); ```
--- 

## Clients orders junction table
This is a many-to-many relation junciton table.
```
create table clients_orders(
    client_id integer references clients(client_id) not null,
    order_id integer references orders(order_id) not null
);```
---

## Employees table
```
create table employees(
    employee_id serial primary key,
    name varchar(50) not null,
    position varchar(50) not null
);```
* name - employee full name.
* position - employee job title(i.e. driver).
---

## Prices
Prices table is a price-list.
```
create table prices(
    price_id serial primary key,
    cargo_type varchar(50) not null unique,
    cost money not null,
    weight integer not null,
    distance integer not null
);```
* cargo_type - price increases if cargo type is fragile or dangerous(i.e. glass, fuel).
* cost - cargo type increasing coefficient(i.e. 1.75 or 175% for fragile cargo).
* weight - cargo weight coefficient; the more weight - the more vehicle needed - the more price.
* distance - cargo distance coefficient.
---

# Database queries
