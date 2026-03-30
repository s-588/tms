-- name: BulkHardDeleteOrders :exec
DELETE FROM orders WHERE order_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteOrders :exec
UPDATE orders SET deleted_at = NOW() WHERE order_id = ANY(sqlc.arg('ids')::int[]);

-- name: GetOrder :one
SELECT 
    o.*,
    c.name as client_name,
    e.name as employee_name,
    t.license_plate as transport_license_plate,
    p.cargo_type as price_cargo_type,
    ns.name as node_start_name,
    ne.name as node_end_name
FROM orders o
LEFT JOIN clients c ON o.client_id = c.client_id
LEFT JOIN employees e ON o.employee_id = e.employee_id
LEFT JOIN transports t ON o.transport_id = t.transport_id
LEFT JOIN prices p ON o.price_id = p.price_id
LEFT JOIN nodes ns ON o.node_id_start = ns.node_id
LEFT JOIN nodes ne ON o.node_id_end = ne.node_id
WHERE o.order_id = sqlc.arg('order_id') AND o.deleted_at IS NULL;

-- name: HardDeleteOrder :exec
DELETE FROM orders WHERE order_id = sqlc.arg('order_id');

-- name: RestoreOrder :exec
UPDATE orders SET deleted_at = NULL WHERE order_id = sqlc.arg('order_id');

-- name: SoftDeleteOrder :exec
UPDATE orders SET deleted_at = NOW() 
WHERE order_id = sqlc.arg('order_id');

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = sqlc.arg('status')::order_status, updated_at = NOW()
WHERE order_id = sqlc.arg('order_id');

-- name: CreateOrder :one
INSERT INTO orders (
    client_id, transport_id, employee_id, grade,
    distance, weight, total_price, price_id, status,
    node_id_start, node_id_end
) VALUES (
    sqlc.arg('client_id')::int,
    sqlc.arg('transport_id')::int,
    sqlc.arg('employee_id')::int,
    sqlc.arg('grade')::smallint,
    sqlc.arg('distance')::double precision,
    sqlc.arg('weight')::int,
    sqlc.arg('total_price')::numeric,
    sqlc.arg('price_id')::int,
    sqlc.arg('status')::order_status,
    sqlc.arg('node_id_start')::int,
    sqlc.arg('node_id_end')::int
)
RETURNING *;

-- name: GetOrders :many
SELECT 
    o.*,
    c.name as client_name,
    e.name as employee_name,
    t.license_plate as transport_license_plate,
    p.cargo_type as price_cargo_type,
    ns.name as node_start_name,
    ne.name as node_end_name,
    (count(*) OVER())/20+1 AS total_count
FROM orders o
LEFT JOIN clients c ON o.client_id = c.client_id
LEFT JOIN employees e ON o.employee_id = e.employee_id
LEFT JOIN transports t ON o.transport_id = t.transport_id
LEFT JOIN prices p ON o.price_id = p.price_id
LEFT JOIN nodes ns ON o.node_id_start = ns.node_id
LEFT JOIN nodes ne ON o.node_id_end = ne.node_id
WHERE o.deleted_at IS NULL
  AND (sqlc.narg('status_filter')::order_status IS NULL OR o.status = sqlc.narg('status_filter')::order_status)
  AND (sqlc.narg('total_price_min')::numeric IS NULL OR o.total_price >= sqlc.narg('total_price_min')::numeric)
  AND (sqlc.narg('total_price_max')::numeric IS NULL OR o.total_price <= sqlc.narg('total_price_max')::numeric)
  AND (sqlc.narg('distance_min')::double precision IS NULL OR o.distance >= sqlc.narg('distance_min')::double precision)
  AND (sqlc.narg('distance_max')::double precision IS NULL OR o.distance <= sqlc.narg('distance_max')::double precision)
  AND (sqlc.narg('weight_min')::int IS NULL OR o.weight >= sqlc.narg('weight_min')::int)
  AND (sqlc.narg('weight_max')::int IS NULL OR o.weight <= sqlc.narg('weight_max')::int)
  AND (sqlc.narg('client_id_filter')::int IS NULL OR o.client_id = sqlc.narg('client_id_filter')::int)
  AND (sqlc.narg('transport_id_filter')::int IS NULL OR o.transport_id = sqlc.narg('transport_id_filter')::int)
  AND (sqlc.narg('employee_id_filter')::int IS NULL OR o.employee_id = sqlc.narg('employee_id_filter')::int)
  AND (sqlc.narg('price_id_filter')::int IS NULL OR o.price_id = sqlc.narg('price_id_filter')::int)
  AND (sqlc.narg('grade_min')::smallint IS NULL OR o.grade >= sqlc.narg('grade_min')::smallint)
  AND (sqlc.narg('grade_max')::smallint IS NULL OR o.grade <= sqlc.narg('grade_max')::smallint)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'order_id' THEN o.order_id::text
            WHEN 'distance' THEN o.distance::text
            WHEN 'weight' THEN o.weight::text
            WHEN 'total_price' THEN o.total_price::text
            WHEN 'status' THEN o.status::text
            WHEN 'grade' THEN o.grade::text
            WHEN 'created_at' THEN o.created_at::text
            WHEN 'updated_at' THEN o.updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'order_id' THEN o.order_id::text
            WHEN 'distance' THEN o.distance::text
            WHEN 'weight' THEN o.weight::text
            WHEN 'total_price' THEN o.total_price::text
            WHEN 'status' THEN o.status::text
            WHEN 'grade' THEN o.grade::text
            WHEN 'created_at' THEN o.created_at::text
            WHEN 'updated_at' THEN o.updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: UpdateOrder :exec
UPDATE orders
SET
    client_id = sqlc.arg('client_id')::int,
    transport_id = sqlc.arg('transport_id')::int,
    employee_id = sqlc.arg('employee_id')::int,
    grade = sqlc.arg('grade')::smallint,
    distance = sqlc.arg('distance')::double precision,  
    weight = sqlc.arg('weight')::int,
    total_price = sqlc.arg('total_price')::numeric,
    price_id = sqlc.arg('price_id')::int,
    status = sqlc.arg('status')::order_status,
    node_id_start = sqlc.arg('node_id_start')::int,
    node_id_end = sqlc.arg('node_id_end')::int,
    updated_at = NOW()
WHERE order_id = sqlc.arg('order_id');