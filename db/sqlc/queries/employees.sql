-- Get paginated employees list
-- name: GetEmployeesPaginated :many
select *,count(*) as total_count   from employees e
where e.deleted_at is null
order by employee_id
limit $1 offset $2;

-- Get single employee by employee_id
-- name: GetEmployeeByemployee_id :one
select * from employees e where employee_id = $1 and e.deleted_at is null;

-- Create new employee
-- name: CreateEmployee :one
insert into employees(name) values ($1) returning *;

-- Delete employee
-- name: DeleteEmployee :exec
delete from employees where employee_id = $1;

-- Update employee name
-- name: UpdateEmployee :exec
update employees set name = $2 where employee_id = $1;
