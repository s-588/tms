-- name: BulkHardDeleteEmployees :exec
DELETE FROM employees WHERE employee_id = ANY(sqlc.arg('ids')::int[]);

-- name: BulkSoftDeleteEmployees :exec
UPDATE employees SET deleted_at = NOW() WHERE employee_id = ANY(sqlc.arg('ids')::int[]);

-- name: CreateEmployee :one
INSERT INTO employees (
    name, status, job_title, hire_date, salary,
    license_issued, license_expiration
) VALUES (
    sqlc.arg('name')::text,
    sqlc.arg('status')::employee_status,
    sqlc.arg('job_title')::employee_job_title,
    sqlc.arg('hire_date')::date,
    sqlc.arg('salary')::numeric,
    sqlc.arg('license_issued')::date,
    sqlc.arg('license_expiration')::date
)
RETURNING *;

-- name: GetEmployee :one
SELECT * FROM employees
WHERE employee_id = sqlc.arg('employee_id') AND deleted_at IS NULL;

-- name: GetEmployees :many
SELECT *,
       (count(*) OVER())/20+1 AS total_count
FROM employees
WHERE deleted_at IS NULL
  AND (sqlc.narg('name_filter')::text IS NULL OR name ILIKE '%' || sqlc.narg('name_filter')::text || '%')
  AND (sqlc.narg('job_title_filter')::employee_job_title IS NULL OR job_title = sqlc.narg('job_title_filter'))
  AND (sqlc.narg('status_filter')::employee_status IS NULL OR status = sqlc.narg('status_filter')::employee_status)
  AND (sqlc.narg('salary_min')::numeric IS NULL OR salary >= sqlc.narg('salary_min')::numeric)
  AND (sqlc.narg('salary_max')::numeric IS NULL OR salary <= sqlc.narg('salary_max')::numeric)
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'ASC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'job_title' THEN job_title::text
            WHEN 'status' THEN status::text
            WHEN 'salary' THEN salary::text
            WHEN 'hire_date' THEN hire_date::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'DESC' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'job_title' THEN job_title::text
            WHEN 'status' THEN status::text
            WHEN 'salary' THEN salary::text
            WHEN 'hire_date' THEN hire_date::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT 20 OFFSET 20 * (sqlc.arg('page')::integer - 1);

-- name: HardDeleteEmployee :exec
DELETE FROM employees WHERE employee_id = sqlc.arg('employee_id');

-- name: RestoreEmployee :exec
UPDATE employees SET deleted_at = NULL WHERE employee_id = sqlc.arg('employee_id');

-- name: SoftDeleteEmployee :exec
UPDATE employees SET deleted_at = NOW() WHERE employee_id = sqlc.arg('employee_id');

-- name: UpdateEmployee :exec
UPDATE employees
SET
    name = sqlc.arg('name')::text,
    status = sqlc.arg('status')::employee_status,
    job_title = sqlc.arg('job_title')::employee_job_title,
    hire_date = sqlc.arg('hire_date')::date,
    salary = sqlc.arg('salary')::numeric,
    license_issued = sqlc.arg('license_issued')::date,
    license_expiration = sqlc.arg('license_expiration')::date,
    updated_at = NOW()
WHERE employee_id = sqlc.arg('employee_id');