# Database schema

## Clients table
```sql
create table clients(
    client_id serial primary key,
    name varchar(50) not null,
    email varchar(254) not null,
    email_verified boolean not null default false,
    phone varchar(25) check (phone ~ '^[\+]?[0-9\-\s()]{7,25}$') not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);```
* name - clients full name
* email - email address for notifications.
* email_verified - indicates if email has been verified
* phone - phone number for extra notifications.
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)
---

## Employees table
```sql
create table employees(
    employee_id serial primary key,
    name varchar(50) not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);```
* name - employee full name.
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)
---

## Fuels table (formerly Fuel types table)
```sql
create table fuels(
    fuel_id serial primary key,
    name varchar(50) not null,
    supplier varchar(50),
    price money not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);```
* name - fuel type name.
* supplier - supplier of the fuel.
* price - fuel price per 1 liter.
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)
---

## Orders table
```sql
create table orders(
    order_id serial primary key,
    distance integer not null,
    weight integer not null,
    total_price money not null,
    status varchar(50) not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
); ```
* distance - distance between cargo loading and end delivery points.
* total_price - total price client will pay.
* status - current order status (e.g., pending, in transit, delivered, etc.)
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)

Total price calculated as: fuel_full_price + (weight * distance * prices_total_coefficient).
Where:
    - fuel_full_price = fuel(price) * orders(distance) * transport(fuel_consumption)
    - prices_total_coefficient = prices(cost) + prices(weight) + prices(distance)
---

## Orders transport junction table
One order can be transported by multiple vehicles.
```sql
create table orders_transport(
    order_id integer references orders(order_id),
    transport_id integer references transports(transport_id)
); ```
--- 

## Clients orders junction table
This is a many-to-many relation junction table.
```sql
create table clients_orders(
    client_id integer references clients(client_id) not null,
    order_id integer references orders(order_id) not null
);```
---

## Transports table
```sql
create table transports(
    transport_id serial primary key,
    employee_id integer references employees(employee_id),
    model varchar(50) not null,
    license_plate varchar(10),
    payload_capacity integer not null,
    fuel_id integer references fuels(fuel_id) not null,
    fuel_consumption integer not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);```
* employee_id - who is currently responsible for this vehicle.
* model - vehicle model name.
* license_plate - vehicle registration plate number.
* payload_capacity - maximum weight vehicle can deliver.
* fuel_id - fuel type.
* fuel_consumption - fuel consumption per 100 km.
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)
---

## Prices table
Prices table is a price-list.
```sql
create table prices(
    price_id serial primary key,
    cargo_type varchar(50) not null unique,
    cost money not null,
    weight integer not null,
    distance integer not null,
    created_at timestamp default now(),
    update_at timestamp default null,
    deleted_at timestamp default null
);```
* cargo_type - price increases if cargo type is fragile or dangerous(i.e. glass, fuel).
* cost - cargo type increasing coefficient(i.e. 1.75 or 175% for fragile cargo).
* weight - cargo weight coefficient; the more weight - the more vehicle needed - the more price.
* distance - cargo distance coefficient.
* created_at - record creation timestamp
* update_at - record last update timestamp
* deleted_at - soft delete timestamp (null if not deleted)
---

# Database queries
