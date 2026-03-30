-- name: BulkHardDeletePrices :exec
DELETE FROM prices WHERE price_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeletePrices :exec
UPDATE prices SET deleted_at = NOW() WHERE price_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreatePrice :one
INSERT INTO prices (cargo_type, weight, distance)
VALUES (
    sqlc.arg('cargo_type')::text,
    sqlc.arg('weight')::decimal(10,2),
    sqlc.arg('distance')::decimal(10,2)
)
RETURNING *;

-- name: GetPrice :one
SELECT * FROM prices
WHERE price_id = sqlc.arg('price_id') AND deleted_at IS NULL;

-- name: GetPriceByUnique :one
SELECT * FROM prices
WHERE cargo_type = sqlc.arg('cargo_type')::text
  AND weight = sqlc.arg('weight')::decimal(10,2)
  AND distance = sqlc.arg('distance')::decimal(10,2)
  AND deleted_at IS NULL;

-- name: GetPrices :many
SELECT *,
       (count(*) OVER())/20+1 AS total_count
FROM prices
WHERE deleted_at IS NULL
  AND (sqlc.narg('cargo_type_filter')::text IS NULL OR cargo_type ILIKE '%' || sqlc.narg('cargo_type_filter')::text || '%')
  AND (sqlc.narg('weight_min')::int IS NULL OR weight >= sqlc.narg('weight_min')::int)
  AND (sqlc.narg('weight_max')::int IS NULL OR weight <= sqlc.narg('weight_max')::int)
  AND (sqlc.narg('distance_min')::int IS NULL OR distance >= sqlc.narg('distance_min')::int)
  AND (sqlc.narg('distance_max')::int IS NULL OR distance <= sqlc.narg('distance_max')::int)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'price_id' THEN price_id::text
            WHEN 'cargo_type' THEN cargo_type
            WHEN 'weight' THEN weight::text
            WHEN 'distance' THEN distance::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'price_id' THEN price_id::text
            WHEN 'cargo_type' THEN cargo_type
            WHEN 'weight' THEN weight::text
            WHEN 'distance' THEN distance::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeletePrice :exec
DELETE FROM prices WHERE price_id = sqlc.arg('price_id');

-- name: RestorePrice :exec
UPDATE prices SET deleted_at = NULL WHERE price_id = sqlc.arg('price_id');

-- name: SoftDeletePrice :exec
UPDATE prices SET deleted_at = NOW() WHERE price_id = sqlc.arg('price_id');

-- name: UpdatePrice :exec
UPDATE prices
SET
    cargo_type = sqlc.arg('cargo_type')::text,
    weight = sqlc.arg('weight')::decimal(10,2),
    distance = sqlc.arg('distance')::decimal(10,2),
    updated_at = NOW()
WHERE price_id = sqlc.arg('price_id');