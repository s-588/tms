-- Get paginated order list
-- name: GetOrderPaginated :many
select *,count(*) as total_count   from orders o
where o.deleted_at is null
order by order_id
limit $1 offset $2;

-- Get single order by order_id
-- name: GetOrderByorder_id :one
select * from orders o where order_id = $1 and o.deleted_at is null;

-- Create new order
-- name: CreateOrder :one
insert into orders (distance,weight,total_price,status)
values ($1,$2,$3,$4)
returning *;

-- Delete order
-- name: DeleteOrder :exec
delete from orders where order_id = $1;

-- Update order fields
-- name: UpdateOrder :exec
update orders
set distance = coalesce($2, distance),
weight = coalesce($3, weight),
total_price = coalesce($4, total_price),
status = coalesce($5, status)
where order_id = $1;

-- Get order's transports
-- name: GetOrderTransports :many
select t.*,count(*) as total_count   from transports t
join orders_transports ot on t.transport_id = ot.transport_id
where ot.order_id = $1 and t.deleted_at is null
order by ot.order_id
limit $2 offset $3;

-- Assign/unassign transports to order (replace all connections)
-- First delete existing assignments
-- name: DeleteOrderTransportAssignments :exec
delete from orders_transports where order_id = $1;
-- Then insert new assignments
-- name: InsertOrderTransportAssignments :exec
insert into orders_transports (order_id, transport_id)
select $1, unnest($2::int[]);
