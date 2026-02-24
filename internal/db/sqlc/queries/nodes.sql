-- name: BulkHardDeleteNodes :exec
DELETE FROM nodes WHERE node_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteNodes :exec
UPDATE nodes SET deleted_at = NOW() WHERE node_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateNode :one
INSERT INTO nodes (
    name, geom
) VALUES (
    sqlc.narg('name')::text,
    sqlc.arg('geom')::point
)
RETURNING *;

-- name: GetNode :one
SELECT * FROM nodes
WHERE node_id = sqlc.arg('node_id') AND deleted_at IS NULL;

-- name: GetNodes :many
SELECT *,
       count(*) OVER() AS total_count
FROM nodes
WHERE deleted_at IS NULL
  AND (sqlc.narg('name_filter')::text IS NULL OR name ILIKE '%' || sqlc.narg('name_filter')::text || '%')
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to')::timestamptz)
  AND (sqlc.narg('updated_from')::timestamptz IS NULL OR updated_at >= sqlc.narg('updated_from')::timestamptz)
  AND (sqlc.narg('updated_to')::timestamptz IS NULL OR updated_at <= sqlc.narg('updated_to')::timestamptz)
ORDER BY
    CASE WHEN sqlc.narg('sort_order')::text = 'ASC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'node_id' THEN node_id::text
            WHEN 'name' THEN name
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.narg('sort_order')::text = 'DESC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'node_id' THEN node_id::text
            WHEN 'name' THEN name
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: HardDeleteNode :exec
DELETE FROM nodes WHERE node_id = sqlc.arg('node_id');

-- name: RestoreNode :exec
UPDATE nodes SET deleted_at = NULL WHERE node_id = sqlc.arg('node_id');

-- name: SoftDeleteNode :exec
UPDATE nodes SET deleted_at = NOW() WHERE node_id = sqlc.arg('node_id');

-- name: UpdateNode :exec
UPDATE nodes
SET
    name = COALESCE(sqlc.narg('name')::text, name),
    geom = COALESCE(sqlc.narg('geom')::point, geom),
    updated_at = NOW()
WHERE node_id = sqlc.arg('node_id');
