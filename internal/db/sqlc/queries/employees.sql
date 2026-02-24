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
    sqlc.arg('job_title')::text,
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
       count(*) OVER() AS total_count
FROM employees
WHERE deleted_at IS NULL
  AND (sqlc.narg('name_filter')::text IS NULL OR name ILIKE '%' || sqlc.narg('name_filter')::text || '%')
  AND (sqlc.narg('job_title_filter')::text IS NULL OR job_title ILIKE '%' || sqlc.narg('job_title_filter')::text || '%')
  AND (sqlc.narg('status_filter')::employee_status IS NULL OR status = sqlc.narg('status_filter')::employee_status)
  AND (sqlc.narg('salary_min')::numeric IS NULL OR salary >= sqlc.narg('salary_min')::numeric)
  AND (sqlc.narg('salary_max')::numeric IS NULL OR salary <= sqlc.narg('salary_max')::numeric)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR created_at <= sqlc.narg('created_to')::timestamptz)
  AND (sqlc.narg('updated_from')::timestamptz IS NULL OR updated_at >= sqlc.narg('updated_from')::timestamptz)
  AND (sqlc.narg('updated_to')::timestamptz IS NULL OR updated_at <= sqlc.narg('updated_to')::timestamptz)
ORDER BY
    CASE WHEN sqlc.narg('sort_order')::text = 'ASC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'job_title' THEN job_title
            WHEN 'status' THEN status::text
            WHEN 'salary' THEN salary::text
            WHEN 'hire_date' THEN hire_date::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END ASC,
    CASE WHEN sqlc.narg('sort_order')::text = 'DESC' THEN
        CASE sqlc.narg('sort_by')::text
            WHEN 'name' THEN name
            WHEN 'job_title' THEN job_title
            WHEN 'status' THEN status::text
            WHEN 'salary' THEN salary::text
            WHEN 'hire_date' THEN hire_date::text
            WHEN 'created_at' THEN created_at::text
            WHEN 'updated_at' THEN updated_at::text
        END
    END DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: HardDeleteEmployee :exec
DELETE FROM employees WHERE employee_id = sqlc.arg('employee_id');

-- name: RestoreEmployee :exec
UPDATE employees SET deleted_at = NULL WHERE employee_id = sqlc.arg('employee_id');

-- name: SoftDeleteEmployee :exec
UPDATE employees SET deleted_at = NOW() WHERE employee_id = sqlc.arg('employee_id');

-- name: UpdateEmployee :exec
UPDATE employees
SET
    name = COALESCE(sqlc.narg('name')::text, name),
    status = COALESCE(sqlc.narg('status')::employee_status, status),
    job_title = COALESCE(sqlc.narg('job_title')::text, job_title),
    hire_date = COALESCE(sqlc.narg('hire_date')::date, hire_date),
    salary = COALESCE(sqlc.narg('salary')::numeric, salary),
    license_issued = COALESCE(sqlc.narg('license_issued')::date, license_issued),
    license_expiration = COALESCE(sqlc.narg('license_expiration')::date, license_expiration),
    updated_at = NOW()
WHERE employee_id = sqlc.arg('employee_id');
