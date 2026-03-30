-- name: BulkHardDeleteTransports :exec
DELETE FROM transports WHERE transport_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteTransports :exec
UPDATE transports SET deleted_at = NOW() WHERE transport_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateTransport :one
INSERT INTO transports (
    model, license_plate, payload_capacity, fuel_consumption,
    inspection_passed, inspection_date
) VALUES (
    sqlc.arg('model')::text,
    sqlc.arg('license_plate')::text,
    sqlc.arg('payload_capacity')::int,
    sqlc.arg('fuel_consumption')::int,
    sqlc.arg('inspection_passed')::boolean,
    sqlc.arg('inspection_date')::date
)
RETURNING *;

-- name: GetTransport :one
SELECT * FROM transports
WHERE transport_id = sqlc.arg('transport_id') AND deleted_at IS NULL;

-- name: GetTransportOrders :many
SELECT o.*,
       count(*) OVER() AS total_count
FROM orders o
WHERE o.transport_id = sqlc.arg('transport_id') AND o.deleted_at IS NULL
ORDER BY o.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetTransports :many
SELECT *,
       (count(*) OVER())/20+1 AS total_count
FROM transports
WHERE deleted_at IS NULL
  AND (sqlc.narg('model_filter')::text IS NULL OR model ILIKE '%' || sqlc.narg('model_filter')::text || '%')
  AND (sqlc.narg('license_plate_filter')::text IS NULL OR license_plate ILIKE '%' || sqlc.narg('license_plate_filter')::text || '%')
  AND (sqlc.narg('payload_capacity_min')::int IS NULL OR payload_capacity >= sqlc.narg('payload_capacity_min')::int)
  AND (sqlc.narg('payload_capacity_max')::int IS NULL OR payload_capacity <= sqlc.narg('payload_capacity_max')::int)
  AND (sqlc.narg('fuel_consumption_min')::int IS NULL OR fuel_consumption >= sqlc.narg('fuel_consumption_min')::int)
  AND (sqlc.narg('fuel_consumption_max')::int IS NULL OR fuel_consumption <= sqlc.narg('fuel_consumption_max')::int)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'model' THEN model
            WHEN 'license_plate' THEN license_plate
            WHEN 'payload_capacity' THEN payload_capacity::text
            WHEN 'fuel_consumption' THEN fuel_consumption::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'model' THEN model
            WHEN 'license_plate' THEN license_plate
            WHEN 'payload_capacity' THEN payload_capacity::text
            WHEN 'fuel_consumption' THEN fuel_consumption::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeleteTransport :exec
DELETE FROM transports WHERE transport_id = sqlc.arg('transport_id');

-- name: RestoreTransport :exec
UPDATE transports SET deleted_at = NULL WHERE transport_id = sqlc.arg('transport_id');

-- name: SoftDeleteTransport :exec
UPDATE transports SET deleted_at = NOW() WHERE transport_id = sqlc.arg('transport_id');

-- name: UpdateTransport :exec
UPDATE transports
SET
    model = sqlc.arg('model')::text,
    license_plate = sqlc.arg('license_plate')::text,
    payload_capacity = sqlc.arg('payload_capacity')::int,
    fuel_consumption = sqlc.arg('fuel_consumption')::int,
    updated_at = NOW()
WHERE transport_id = sqlc.arg('transport_id');