-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_clients_deleted_at ON clients(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_fuels_deleted_at ON fuels(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_prices_deleted_at ON prices(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_transports_deleted_at ON transports(deleted_at) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_transports_employee_id ON transports(employee_id);
CREATE INDEX idx_transports_fuel_id ON transports(fuel_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_clients_email_token ON clients(email_token) WHERE email_token IS NOT NULL;

-- For email uniqueness
CREATE UNIQUE INDEX idx_clients_email_unique ON clients(email) WHERE deleted_at IS NULL;

CREATE INDEX idx_orders_status ON orders(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_clients_deleted_at;
DROP INDEX IF EXISTS idx_employees_deleted_at;
DROP INDEX IF EXISTS idx_fuels_deleted_at;
DROP INDEX IF EXISTS idx_orders_deleted_at;
DROP INDEX IF EXISTS idx_prices_deleted_at;
DROP INDEX IF EXISTS idx_transports_deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_clients_orders_client_id;
DROP INDEX IF EXISTS idx_clients_orders_order_id;
DROP INDEX IF EXISTS idx_orders_transports_order_id;
DROP INDEX IF EXISTS idx_orders_transports_transport_id;
DROP INDEX IF EXISTS idx_transports_employee_id;
DROP INDEX IF EXISTS idx_transports_fuel_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_clients_email_token;
DROP INDEX IF EXISTS idx_clients_email_unique;
DROP INDEX IF EXISTS idx_orders_status;
-- +goose StatementEnd
