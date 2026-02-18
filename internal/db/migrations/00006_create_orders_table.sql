-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
drop type if exists order_status;
-- +goose StatementEnd
