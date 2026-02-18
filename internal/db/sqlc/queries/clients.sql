-- Get paginated client list
-- name: GetClientsPaginated :many
SELECT *,count(*) over() as total_count FROM clients
WHERE deleted_at IS NULL
ORDER BY client_id
LIMIT $1 OFFSET $2 ;

-- Get single client by client_id
-- name: GetClientByclient_id :one
SELECT * FROM clients WHERE client_id = $1 AND deleted_at IS NULL;

-- Create new client
-- name: CreateClient :one
INSERT INTO clients (name, email, phone, email_verified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Soft delete client (set deleted_at)
-- name: DeleteClient :exec
delete from clients WHERE client_id = $1;

-- Update client fields
-- name: UpdateClient :exec
UPDATE clients
SET
  name = COALESCE($2, name),
  email = COALESCE($3, email),
  email_verified = CASE WHEN $3 IS NOT NULL THEN false ELSE email_verified END,
  phone = COALESCE($4, phone)
WHERE client_id = $1;

-- Verify client email by token
-- name: VerifyClientEmail :exec
UPDATE clients SET email_verified = true WHERE email_token = $1;

-- Get client's orders
-- name: GetClientOrders :many
SELECT o.*,count(*) over() as total_count  FROM orders o
JOIN clients_orders co ON o.order_id = co.order_id
WHERE co.client_id = $1 and o.deleted_at is null
ORDER BY co.client_id
LIMIT $2 OFFSET $3;

-- Assign/unassign client to orders (replace all connections)
-- First delete existing assignments
-- name: DeleteClientOrderAssignments :exec
DELETE FROM clients_orders WHERE client_id = $1;
-- Then insert new assignments
-- name: InsertClientOrderAssignments :exec
INSERT INTO clients_orders (client_id, order_id)
SELECT $1, unnest($2::int[]);
