-- name: BulkHardDeleteNodes :exec
DELETE FROM nodes WHERE node_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteNodes :exec
UPDATE nodes SET deleted_at = NOW() WHERE node_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateNode :one
INSERT INTO nodes (
    address, name, geom
) VALUES (
    sqlc.arg('address')::text,
    sqlc.narg('name')::text,
    sqlc.arg('geom')::geography
)
RETURNING *;

-- name: GetNode :one
SELECT  node_id,
    name,
 ST_X(geom::geometry) AS x,
    ST_Y(geom::geometry) AS y,
    created_at,
    updated_at,
    deleted_at
FROM nodes
WHERE node_id = sqlc.arg('node_id') AND deleted_at IS NULL;

-- name: GetNodes :many
SELECT  node_id,
 name,
    ST_X(geom::geometry)::double precision AS x,
    ST_Y(geom::geometry)::double precision AS y,
    created_at,
    updated_at,
    deleted_at,
    (count(*) OVER())/20+1 AS total_count
FROM nodes
WHERE deleted_at IS NULL
  AND (sqlc.narg('name_filter')::text IS NULL OR name ILIKE '%' || sqlc.narg('name_filter')::text || '%')
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'node_id' THEN node_id::text
            WHEN 'name' THEN name
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'node_id' THEN node_id::text
            WHEN 'name' THEN name
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeleteNode :exec
DELETE FROM nodes WHERE node_id = sqlc.arg('node_id');

-- name: RestoreNode :exec
UPDATE nodes SET deleted_at = NULL WHERE node_id = sqlc.arg('node_id');

-- name: SoftDeleteNode :exec
UPDATE nodes SET deleted_at = NOW() WHERE node_id = sqlc.arg('node_id');

-- name: UpdateNode :exec
UPDATE nodes
SET
    address = sqlc.arg('address')::text,
    name = sqlc.narg('name')::text,
    geom = sqlc.arg('geom')::geography,
    updated_at = NOW()
WHERE node_id = sqlc.arg('node_id');

-- name: GetDistanceBetweenNodes :one
SELECT ST_Distance(n1.geom, n2.geom)::float
		FROM nodes n1, nodes n2
		WHERE n1.node_id = $1 AND n2.node_id = $2;