-- Get paginated prices list
-- name: GetPricesPaginated :many
select *,count(*) over() as total_count   from prices p
where p.deleted_at is null
order by price_id
limit $1 offset $2;

-- Get single price by price_id
-- name: GetPriceByprice_id :one
select * from prices p where price_id = $1 and p.deleted_at is null;

-- Create new price
-- name: CreatePrice :one
insert into prices(cargo_type, cost, weight, distance)
values ($1, $2, $3, $4)
returning *;

-- Delete price
-- name: DeletePrice :exec
delete from prices where price_id = $1;

-- Update price fields
-- name: UpdatePrice :exec
update prices
set
cargo_type = coalesce($2, cargo_type),
cost = coalesce($3 ,cost),
weight = coalesce($4, weight),
distance = coalesce($5 , distance)
where price_id = $1;
