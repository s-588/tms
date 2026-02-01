-- Get paginated transports list
-- name: GetTransportsPaginated :many
select *,count(*) over() as total_count   from transports t
where t.deleted_at is null
order by transport_id
limit $1 offset $2;

-- Get single transport by transport_id
-- name: GetTransportBytransport_id :one
select * from transports t where transport_id = $1 and t.deleted_at is null;

-- Create new transport
-- name: CreateTransport :one
insert into transports(employee_id, model, license_plate, payload_capacity, fuel_id, fuel_consumption)
values ($1,$2,$3,$4,$5,$6)
returning *;

-- Delete transport
-- name: DeleteTransport :exec
delete from transports where transport_id = $1;

-- Update transport fields
-- name: UpdateTransport :exec
update transports
set
employee_id = coalesce($2, employee_id),
model = coalesce($3, model),
license_plate = coalesce($4, license_plate),
payload_capacity = coalesce($5, payload_capacity),
fuel_id = coalesce($6, fuel_id),
fuel_consumption = coalesce($7, fuel_consumption)
where transport_id = $1;
