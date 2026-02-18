-- +goose Up
-- +goose StatementBegin
CREATE TYPE employee_status AS ENUM (
    'available',
    'assigned',
    'unavailable'
);

create type employee_job_title as enum(
'driver',
'dispatcher',
'mechanic',
'logistics_manager'
);

CREATE TABLE employees (
    employee_id         SERIAL PRIMARY KEY,
    name                VARCHAR(50) NOT NULL,
    status              employee_status NOT NULL DEFAULT 'available',
    job_title           employee_job_title NOT NULL,
    hire_date           DATE DEFAULT CURRENT_DATE,
    salary              NUMERIC(10,2) NOT NULL
                            CHECK (salary > 0),
    license_issued      DATE NOT NULL,
    license_expiration  DATE NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CHECK (license_expiration > license_issued)
);

CREATE INDEX idx_employees_status
    ON employees (status);

CREATE INDEX idx_employees_job_title
    ON employees (job_title);

CREATE INDEX idx_employees_deleted_at
    ON employees (deleted_at)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
drop type if exists employee_status;
drop type if exists employee_job_title;
-- +goose StatementEnd
