-- name: BulkHardDeleteOrders :exec
DELETE FROM orders WHERE order_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteOrders :exec
UPDATE orders SET deleted_at = NOW() WHERE order_id = ANY(sqlc.arg('ids')::int[]);

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
    sqlc.arg('distance')::int,
    sqlc.arg('weight')::int,
    sqlc.arg('total_price')::numeric,
    sqlc.arg('price_id')::int,
    sqlc.arg('status')::order_status,
    sqlc.arg('node_id_start')::int,
    sqlc.arg('node_id_end')::int
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE order_id = sqlc.arg('order_id') AND deleted_at IS NULL;

-- name: GetOrders :many
SELECT *,
       count(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
  AND (sqlc.narg('status_filter')::order_status IS NULL OR status = sqlc.narg('status_filter')::order_status)
  AND (sqlc.narg('total_price_min')::numeric IS NULL OR total_price >= sqlc.narg('total_price_min')::numeric)
  AND (sqlc.narg('total_price_max')::numeric IS NULL OR total_price <= sqlc.narg('total_price_max')::numeric)
  AND (sqlc.narg('distance_min')::int IS NULL OR distance >= sqlc.narg('distance_min')::int)
  AND (sqlc.narg('distance_max')::int IS NULL OR distance <= sqlc.narg('distance_max')::int)
  AND (sqlc.narg('weight_min')::int IS NULL OR weight >= sqlc.narg('weight_min')::int)
  AND (sqlc.narg('weight_max')::int IS NULL OR weight <= sqlc.narg('weight_max')::int)
  AND (sqlc.narg('client_id_filter')::int IS NULL OR client_id = sqlc.narg('client_id_filter')::int)
  AND (sqlc.narg('transport_id_filter')::int IS NULL OR transport_id = sqlc.narg('transport_id_filter')::int)
  AND (sqlc.narg('employee_id_filter')::int IS NULL OR employee_id = sqlc.narg('employee_id_filter')::int)
  AND (sqlc.narg('price_id_filter')::int IS NULL OR price_id = sqlc.narg('price_id_filter')::int)
  AND (sqlc.narg('grade_min')::smallint IS NULL OR grade >= sqlc.narg('grade_min')::smallint)
  AND (sqlc.narg('grade_max')::smallint IS NULL OR grade <= sqlc.narg('grade_max')::smallint)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to')::timestamptz)
  AND (sqlc.narg('updated_from')::timestamptz IS NULL OR updated_at >= sqlc.narg('updated_from')::timestamptz)
  AND (sqlc.narg('updated_to')::timestamptz IS NULL OR updated_at <= sqlc.narg('updated_to')::timestamptz)
ORDER BY
    CASE WHEN sqlc.narg('sort_order')::text = 'ASC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'order_id' THEN order_id::text
            WHEN 'distance' THEN distance::text
            WHEN 'weight' THEN weight::text
            WHEN 'total_price' THEN total_price::text
            WHEN 'status' THEN status::text
            WHEN 'grade' THEN grade::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.narg('sort_order')::text = 'DESC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'order_id' THEN order_id::text
            WHEN 'distance' THEN distance::text
            WHEN 'weight' THEN weight::text
            WHEN 'total_price' THEN total_price::text
            WHEN 'status' THEN status::text
            WHEN 'grade' THEN grade::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: HardDeleteOrder :exec
DELETE FROM orders WHERE order_id = sqlc.arg('order_id');

-- name: RestoreOrder :exec
UPDATE orders SET deleted_at = NULL WHERE order_id = sqlc.arg('order_id');

-- name: SoftDeleteOrder :exec
UPDATE orders SET deleted_at = NOW() WHERE order_id = sqlc.arg('order_id');

-- name: UpdateOrder :exec
UPDATE orders
SET
    client_id = COALESCE(sqlc.narg('client_id')::int, client_id),
    transport_id = COALESCE(sqlc.narg('transport_id')::int, transport_id),
    employee_id = COALESCE(sqlc.narg('employee_id')::int, employee_id),
    grade = COALESCE(sqlc.narg('grade')::smallint, grade),
    distance = COALESCE(sqlc.narg('distance')::int, distance),
    weight = COALESCE(sqlc.narg('weight')::int, weight),
    total_price = COALESCE(sqlc.narg('total_price')::numeric, total_price),
    price_id = COALESCE(sqlc.narg('price_id')::int, price_id),
    status = COALESCE(sqlc.narg('status')::order_status, status),
    node_id_start = COALESCE(sqlc.narg('node_id_start')::int, node_id_start),
    node_id_end = COALESCE(sqlc.narg('node_id_end')::int, node_id_end),
    updated_at = NOW()
WHERE order_id = sqlc.arg('order_id');

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = sqlc.arg('status')::order_status, updated_at = NOW()
WHERE order_id = sqlc.arg('order_id');
