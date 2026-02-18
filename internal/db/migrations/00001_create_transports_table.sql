-- +goose Up
-- +goose StatementBegin
CREATE TABLE transports (
    transport_id        SERIAL PRIMARY KEY,
    model               VARCHAR(50) NOT NULL,
    license_plate       VARCHAR(10) UNIQUE,
    payload_capacity    INTEGER NOT NULL,
        fuel_consumption    INTEGER NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_transports_deleted_at
    ON transports (deleted_at)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
create table insurances(
    insurance_id serial primary key,
    transport_id integer not null references transports(transport_id) on delete cascade,
    insurance_date DATE NOT NULL,
    insurance_expiration DATE NOT NULL,
    payment numeric(10,2) not null,
    coverage numeric(10,2) not null,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);
create index idx_insurances_delete_at on insurances(deleted_at) where
deleted_at is null;
-- +goose StatementEnd

-- +goose StatementBegin
create type inspection_status as enum('ready','repair','overdue');
create table inspections(
    inspection_id serial primary key,
    transport_id integer not null references transports(transport_id) on delete cascade,
    inspection_date date not null,
    inspection_expiration date not null,
    status inspection_status not null default 'overdue',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);
create index idx_ispections_status on inspections(status);
create index idx_ispections_deleted_at on inspections(deleted_at) where deleted_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists insurances;
drop table if exists inspections;
drop type if exists inspection_status;
DROP TABLE IF EXISTS transports;
-- +goose StatementEnd
