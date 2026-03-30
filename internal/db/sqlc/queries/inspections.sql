-- name: BulkHardDeleteInspections :exec
DELETE FROM inspections WHERE inspection_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteInspections :exec
UPDATE inspections SET deleted_at = NOW() WHERE inspection_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateInspection :one
INSERT INTO inspections (
    transport_id, inspection_date, inspection_expiration, status
) VALUES (
    sqlc.arg('transport_id')::int,
    sqlc.arg('inspection_date')::date,
    sqlc.arg('inspection_expiration')::date,
    sqlc.arg('status')::inspection_status
)
RETURNING *;

-- name: GetInspection :one
SELECT * FROM inspections
WHERE inspection_id = sqlc.arg('inspection_id') AND deleted_at IS NULL;

-- name: GetInspectionsByTransport :many
SELECT * FROM inspections
WHERE transport_id = sqlc.arg('transport_id') AND deleted_at IS NULL
ORDER BY inspection_date DESC;

-- name: GetInspections :many
SELECT *,
       (count(*) OVER())/20+1 AS total_count
FROM inspections
WHERE deleted_at IS NULL
  AND (sqlc.narg('transport_id_filter')::int IS NULL OR transport_id = sqlc.narg('transport_id_filter')::int)
  AND (sqlc.narg('status_filter')::inspection_status IS NULL OR status = sqlc.narg('status_filter')::inspection_status)
  AND (sqlc.narg('inspection_date_from')::date IS NULL OR inspection_date >= sqlc.narg('inspection_date_from')::date)
  AND (sqlc.narg('inspection_date_to')::date IS NULL OR inspection_date <= sqlc.narg('inspection_date_to')::date)
  AND (sqlc.narg('inspection_expiration_from')::date IS NULL OR inspection_expiration >= sqlc.narg('inspection_expiration_from')::date)
  AND (sqlc.narg('inspection_expiration_to')::date IS NULL OR inspection_expiration <= sqlc.narg('inspection_expiration_to')::date)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'inspection_id' THEN inspection_id::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'status' THEN status::text
            WHEN 'inspection_date' THEN inspection_date::text
            WHEN 'inspection_expiration' THEN inspection_expiration::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'inspection_id' THEN inspection_id::text
            WHEN 'transport_id' THEN transport_id::text
            WHEN 'status' THEN status::text
            WHEN 'inspection_date' THEN inspection_date::text
            WHEN 'inspection_expiration' THEN inspection_expiration::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeleteInspection :exec
DELETE FROM inspections WHERE inspection_id = sqlc.arg('inspection_id');

-- name: RestoreInspection :exec
UPDATE inspections SET deleted_at = NULL WHERE inspection_id = sqlc.arg('inspection_id');

-- name: SoftDeleteInspection :exec
UPDATE inspections SET deleted_at = NOW() WHERE inspection_id = sqlc.arg('inspection_id');

-- name: UpdateInspection :exec
UPDATE inspections
SET
    transport_id = sqlc.arg('transport_id')::int,
    inspection_date = sqlc.arg('inspection_date')::date,
    inspection_expiration = sqlc.arg('inspection_expiration')::date,
    status = sqlc.arg('status')::inspection_status,
    updated_at = NOW()
WHERE inspection_id = sqlc.arg('inspection_id');