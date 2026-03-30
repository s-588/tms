-- name: BulkHardDeleteInsurances :exec
DELETE FROM insurances WHERE insurance_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteInsurances :exec
UPDATE insurances SET deleted_at = NOW() WHERE insurance_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateInsurance :one
INSERT INTO insurances (
    transport_id, insurance_date, insurance_expiration, payment, coverage
) VALUES (
    sqlc.arg('transport_id')::int,
    sqlc.arg('insurance_date')::date,
    sqlc.arg('insurance_expiration')::date,
    sqlc.arg('payment')::numeric,
    sqlc.arg('coverage')::numeric
)
RETURNING *;

-- name: GetInsurance :one
SELECT * FROM insurances
WHERE insurance_id = sqlc.arg('insurance_id') AND deleted_at IS NULL;

-- name: GetInsuranceByTransport :one
SELECT * FROM insurances
WHERE transport_id = sqlc.arg('transport_id') AND deleted_at IS NULL
ORDER BY insurance_expiration DESC
LIMIT 1;

-- name: GetInsurances :many
SELECT *,
       (count(*) OVER())/20+1 AS total_count
FROM insurances
WHERE deleted_at IS NULL
  AND (sqlc.narg('transport_id_filter')::int IS NULL OR transport_id = sqlc.narg('transport_id_filter')::int)
  AND (sqlc.narg('insurance_date_from')::date IS NULL OR insurance_date >= sqlc.narg('insurance_date_from')::date)
  AND (sqlc.narg('insurance_date_to')::date IS NULL OR insurance_date <= sqlc.narg('insurance_date_to')::date)
  AND (sqlc.narg('insurance_expiration_from')::date IS NULL OR insurance_expiration >= sqlc.narg('insurance_expiration_from')::date)
  AND (sqlc.narg('insurance_expiration_to')::date IS NULL OR insurance_expiration <= sqlc.narg('insurance_expiration_to')::date)
  AND (sqlc.narg('payment_min')::numeric IS NULL OR payment >= sqlc.narg('payment_min')::numeric)
  AND (sqlc.narg('payment_max')::numeric IS NULL OR payment <= sqlc.narg('payment_max')::numeric)
  AND (sqlc.narg('coverage_min')::numeric IS NULL OR coverage >= sqlc.narg('coverage_min')::numeric)
  AND (sqlc.narg('coverage_max')::numeric IS NULL OR coverage <= sqlc.narg('coverage_max')::numeric)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'insurance_id' THEN insurance_id::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'insurance_date' THEN insurance_date::text
            WHEN 'insurance_expiration' THEN insurance_expiration::text
            WHEN 'payment' THEN payment::text
            WHEN 'coverage' THEN coverage::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'insurance_id' THEN insurance_id::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'insurance_date' THEN insurance_date::text
            WHEN 'insurance_expiration' THEN insurance_expiration::text
            WHEN 'payment' THEN payment::text
            WHEN 'coverage' THEN coverage::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeleteInsurance :exec
DELETE FROM insurances WHERE insurance_id = sqlc.arg('insurance_id');

-- name: RestoreInsurance :exec
UPDATE insurances SET deleted_at = NULL WHERE insurance_id = sqlc.arg('insurance_id');

-- name: SoftDeleteInsurance :exec
UPDATE insurances SET deleted_at = NOW() WHERE insurance_id = sqlc.arg('insurance_id');

-- name: UpdateInsurance :exec
UPDATE insurances
SET
    transport_id = sqlc.arg('transport_id')::int,
    insurance_date = sqlc.arg('insurance_date')::date,
    insurance_expiration = sqlc.arg('insurance_expiration')::date,
    payment = sqlc.arg('payment')::numeric,
    coverage = sqlc.arg('coverage')::numeric,
    updated_at = NOW()
WHERE insurance_id = sqlc.arg('insurance_id');