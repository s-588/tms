-- Get paginated fuel types list
-- name: GetFuelsPaginated :many
select *,count(*) over() as total_count   from fuels f
where f.deleted_at is null
order by fuel_id
limit $1 offset $2;

-- Get single fuel type by fuel_id
-- name: GetFuelByfuel_id :one
select * from fuels f where fuel_id = $1 and f.deleted_at is null;

-- Create new fuel type
-- name: CreateFuel :one
insert into fuels(name,supplier,price)
values ($1,$2,$3)
returning *;

-- Delete fuel type
-- name: DeleteFuel :exec
DELETE from fuels where fuel_id = $1;

-- Update fuel type fields
-- name: UpdateFuel :exec
Update fuels
set
name = coalesce($2, name),
supplier = coalesce($3,supplier),
price = coalesce($4, price)
where fuel_id = $1;
