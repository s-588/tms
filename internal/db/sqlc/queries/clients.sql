-- name: BulkHardDeleteClients :exec
DELETE FROM clients WHERE client_id = ANY($1::int[]);

-- name: BulkSoftDeleteClients :exec
UPDATE clients SET deleted_at = NOW() WHERE client_id = ANY($1::int[]);

-- name: CreateClient :one
INSERT INTO clients (
    name,
    email,
    phone
) VALUES (
    sqlc.arg('name')::text,
    sqlc.arg('email')::text,
    sqlc.arg('phone')::text
)
RETURNING *;

-- name: GetClient :one
SELECT * FROM clients
WHERE client_id = sqlc.arg('client_id') AND deleted_at IS NULL;

-- name: GetClientByEmail :one
SELECT * FROM clients
WHERE email = sqlc.arg('email') AND deleted_at IS NULL;

-- name: GetClientOrders :many
SELECT o.*,
       count(*) OVER() AS total_count
FROM orders o
WHERE o.client_id = sqlc.arg('client_id') AND o.deleted_at IS NULL
ORDER BY o.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetClients :many
SELECT *,
       (count(*) OVER())/20 AS total_count
FROM clients
WHERE deleted_at IS NULL
  AND (sqlc.narg('name_filter')::text IS NULL OR name ILIKE '%' || sqlc.narg('name_filter')::text || '%')
  AND (sqlc.narg('email_filter')::text IS NULL OR email ILIKE '%' || sqlc.narg('email_filter')::text || '%')
  AND (sqlc.narg('phone_filter')::text IS NULL OR phone ILIKE '%' || sqlc.narg('phone_filter')::text || '%')
  AND (sqlc.narg('email_verified_filter')::boolean IS NULL OR email_verified = sqlc.narg('email_verified_filter')::boolean)
  AND (sqlc.narg('score_min_filter')::smallint IS NULL OR score >= sqlc.narg('score_min_filter')::smallint)
  AND (sqlc.narg('score_max_filter')::smallint IS NULL OR score <= sqlc.narg('score_max_filter')::smallint)
  AND (sqlc.narg('created_from_filter')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from_filter')::timestamptz)
  AND (sqlc.narg('created_to_filter')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to_filter')::timestamptz)
  AND (sqlc.narg('updated_from_filter')::timestamptz IS NULL OR updated_at >= sqlc.narg('updated_from_filter')::timestamptz)
  AND (sqlc.narg('updated_to_filter')::timestamptz IS NULL OR updated_at <= sqlc.narg('updated_to_filter')::timestamptz)
ORDER BY
    CASE WHEN sqlc.narg('sort_order')::text = 'ASC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'email' THEN email
            WHEN 'phone' THEN phone
            WHEN 'email_verified' THEN email_verified::text
            WHEN 'score' THEN score::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.narg('sort_order')::text = 'DESC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'email' THEN email
            WHEN 'phone' THEN phone
            WHEN 'email_verified' THEN email_verified::text
            WHEN 'score' THEN score::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20*(sqlc.arg('page')::integer-1);

-- name: HardDeleteClient :exec
DELETE FROM clients WHERE client_id = $1;

-- name: RestoreClient :exec
UPDATE clients SET deleted_at = NULL WHERE client_id = $1;

-- name: SetEmailVerificationToken :exec
UPDATE clients
SET
    email_token = sqlc.arg('email_token')::text,
    email_token_expiration = sqlc.arg('email_token_expiration')::timestamptz
WHERE client_id = sqlc.arg('client_id');

-- name: SoftDeleteClient :exec
UPDATE clients SET deleted_at = NOW() WHERE client_id = $1;

-- name: UpdateClient :exec
UPDATE clients
SET
    name = COALESCE(sqlc.narg('name')::text, name),
    email = COALESCE(sqlc.narg('email')::text, email),
    phone = COALESCE(sqlc.narg('phone')::text, phone),
    updated_at = NOW()
WHERE client_id = sqlc.arg('client_id');

-- name: VerifyClientEmail :exec
UPDATE clients
SET
    email_verified = true,
    email_token = NULL,
    email_token_expiration = NULL
WHERE email_token = sqlc.arg('email_token');
