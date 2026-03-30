-- name: ListClients :many
SELECT
    client_id,
    name
FROM
    clients
WHERE
    deleted_at IS NULL
ORDER BY
    name;

-- name: ListEmployees :many
SELECT
    employee_id,
    name,
    status,
    job_title
FROM
    employees
WHERE
    deleted_at IS NULL
ORDER BY
    name;

-- name: ListFreeDrivers :many
SELECT
    employee_id,
    name,
    status,
    job_title
FROM
    employees
WHERE
    deleted_at IS NULL
    and status = 'available'
    and job_title = 'driver'
ORDER BY
    name;

-- name: ListFreeTransports :many
SELECT
    t.transport_id,
    model,
    license_plate,
    payload_capacity
FROM
    transports t
    left join orders o on t.transport_id = o.transport_id
    inner join insurances i on t.transport_id = i.transport_id
WHERE
    t.deleted_at IS NULL
    and o.transport_id is null
    and i.insurance_expiration > NOW ()
ORDER BY
    license_plate;

-- name: ListTransports :many
SELECT
    transport_id,
    model,
    license_plate
FROM
    transports
WHERE
    deleted_at IS NULL
ORDER BY
    license_plate;

-- name: ListPrices :many
SELECT
    price_id,
    cargo_type,
    weight,
    distance
FROM
    prices
WHERE
    deleted_at IS NULL
ORDER BY
    cargo_type,
    weight,
    distance;

-- name: ListNodes :many
SELECT
    node_id,
    address
FROM
    nodes
WHERE
    deleted_at IS NULL
ORDER BY
    address;