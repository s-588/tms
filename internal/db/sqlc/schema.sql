-- +goose Up (combined from all migration files)

-- 00001_create_transports_table.sql
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

CREATE TABLE insurances(
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

CREATE INDEX idx_insurances_delete_at ON insurances(deleted_at) WHERE deleted_at IS NULL;

CREATE TYPE inspection_status AS ENUM('ready','repair','overdue');

CREATE TABLE inspections(
    inspection_id serial primary key,
    transport_id integer not null references transports(transport_id) on delete cascade,
    inspection_date date not null,
    inspection_expiration date not null,
    status inspection_status not null default 'overdue',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_ispections_status ON inspections(status);
CREATE INDEX idx_ispections_deleted_at ON inspections(deleted_at) WHERE deleted_at IS NULL;

-- 00002_create_clients_table.sql
CREATE TABLE clients (
    client_id               SERIAL PRIMARY KEY,
    name                    VARCHAR(50)  NOT NULL,
    email                   VARCHAR(254) NOT NULL UNIQUE,
    email_verified          BOOLEAN      NOT NULL DEFAULT FALSE,
    email_token             VARCHAR(128) UNIQUE,
    email_token_expiration  TIMESTAMPTZ,
    phone                   VARCHAR(25)  NOT NULL UNIQUE
                                CHECK (phone ~ '^[\+]?[0-9\-\s()]{7,25}$'),
    score                   SMALLINT     NOT NULL DEFAULT 0
                                CHECK (score BETWEEN 0 AND 100),
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ
);

CREATE INDEX idx_clients_deleted_at
    ON clients (deleted_at)
    WHERE deleted_at IS NULL;

-- 00003_create_prices_table.sql
CREATE TABLE prices (
    price_id    SERIAL PRIMARY KEY,
    cargo_type  VARCHAR(50) NOT NULL,
    weight      INTEGER NOT NULL,
    distance    INTEGER NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    UNIQUE (cargo_type, weight, distance)
);

CREATE INDEX idx_prices_deleted_at
    ON prices (deleted_at)
    WHERE deleted_at IS NULL;

-- 00004_create_employees_table.sql
CREATE TYPE employee_status AS ENUM (
    'available',
    'assigned',
    'unavailable'
);

CREATE TYPE employee_job_title AS ENUM(
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

-- 00005_create_nodes_table.sql
CREATE TABLE nodes(
    node_id serial primary key,
    name varchar(50),
    geom point not null,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_nodes_geom ON nodes USING spgist(geom);
CREATE INDEX idx_nodes_name ON nodes(name);
CREATE INDEX idx_nodes_deleted_at ON nodes(deleted_at) WHERE deleted_at IS NULL;

-- 00006_create_orders_table.sql
CREATE TYPE order_status AS ENUM (
    'pending',
    'assigned',
    'in_progress',
    'completed',
    'cancelled'
);

CREATE TABLE orders (
    order_id     SERIAL PRIMARY KEY,
    client_id    INTEGER NOT NULL
                    REFERENCES clients (client_id)
                    ON DELETE RESTRICT ,
    transport_id INTEGER not null
                    REFERENCES transports (transport_id)
                    ON DELETE restrict ,
    employee_id  INTEGER not null
                    REFERENCES employees (employee_id)
                    ON DELETE restrict,
    price_id     INTEGER NOT NULL
                    REFERENCES prices (price_id)
                    ON DELETE RESTRICT,
    grade        SMALLINT NOT NULL
                    CHECK (grade BETWEEN 0 AND 100),
    distance     INTEGER NOT NULL
                    CHECK (distance > 0),
    weight       INTEGER NOT NULL
                    CHECK (weight > 0),
    total_price  NUMERIC(10,2) NOT NULL
                    CHECK (total_price > 0),
    status       order_status NOT NULL DEFAULT 'pending',
    node_id_start integer references nodes(node_id) on delete restrict,
    node_id_end integer references nodes(node_id) on delete restrict,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_orders_client_id
    ON orders (client_id);

CREATE INDEX idx_orders_transport_id
    ON orders (transport_id);

CREATE INDEX idx_orders_employee_id
    ON orders (employee_id);

CREATE INDEX idx_orders_node_id_start
    ON orders (node_id_start);

CREATE INDEX idx_orders_node_id_end
    ON orders (node_id_end);

CREATE INDEX idx_orders_status
    ON orders (status);

CREATE INDEX idx_orders_deleted_at
    ON orders (deleted_at)
    WHERE deleted_at IS NULL;

-- 00007_util_functions.sql
CREATE FUNCTION set_updated_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 00008_transports_functions.sql
CREATE FUNCTION set_transports_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE transports
    SET deleted_at = now()
    WHERE transport_id = OLD.transport_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION set_insurances_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE insurances
    SET deleted_at = now()
    WHERE insurance_id = OLD.insurance_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION set_inspections_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE inspections
    SET deleted_at = now()
    WHERE inspection_id = OLD.inspection_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- 00009_clients_functions.sql
CREATE FUNCTION set_clients_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE clients
    SET deleted_at = now()
    WHERE client_id = OLD.client_id;
    RETURN null;
END;
$$ LANGUAGE plpgsql;

-- 00010_prices_functions.sql
CREATE FUNCTION set_prices_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE prices
    SET deleted_at = now()
    WHERE price_id = OLD.price_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- 00011_employees_functions.sql
CREATE FUNCTION set_employees_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE employees
    SET deleted_at = now()
    WHERE employee_id = OLD.employee_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- 00012_nodes_functions.sql
CREATE FUNCTION set_nodes_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE nodes
    SET deleted_at = now()
    WHERE node_id = OLD.node_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- 00013_orders_functions.sql
CREATE FUNCTION set_orders_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET deleted_at = now()
    WHERE order_id = OLD.order_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- 00014_transports_triggers.sql
CREATE TRIGGER trigger_transports_updated_at
    BEFORE UPDATE ON transports
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_transports_deleted_at
    BEFORE DELETE ON transports
    FOR EACH ROW
    EXECUTE FUNCTION set_transports_deleted_at_trigger_func();

-- 00015_clients_triggers.sql
CREATE TRIGGER trigger_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_clients_deleted_at
    BEFORE DELETE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION set_clients_deleted_at_trigger_func();

-- 00016_prices_triggers.sql
CREATE TRIGGER trigger_prices_updated_at
    BEFORE UPDATE ON prices
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_prices_deleted_at
    BEFORE DELETE ON prices
    FOR EACH ROW
    EXECUTE FUNCTION set_prices_deleted_at_trigger_func();

-- 00017_employees_triggers.sql
CREATE TRIGGER trigger_employees_updated_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_employees_deleted_at
    BEFORE DELETE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION set_employees_deleted_at_trigger_func();

-- 00018_nodes_triggers.sql
CREATE TRIGGER trigger_nodes_updated_at
    BEFORE UPDATE ON nodes
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_nodes_deleted_at
    BEFORE DELETE ON nodes
    FOR EACH ROW
    EXECUTE FUNCTION set_nodes_deleted_at_trigger_func();

-- 00019_orders_triggers.sql
CREATE TRIGGER trigger_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_trigger_func();

CREATE TRIGGER trigger_orders_deleted_at
    BEFORE DELETE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_orders_deleted_at_trigger_func();

-- 00020_orders_calculate_functions.sql (placeholder, no permanent objects)
