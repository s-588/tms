-- +goose Up

-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd

