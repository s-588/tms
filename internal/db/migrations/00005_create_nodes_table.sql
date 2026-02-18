-- +goose Up
-- +goose StatementBegin
create table nodes(
    node_id serial primary key,
    name varchar(50),
    geom point not null,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ
);
create index idx_nodes_geom on nodes using spgist(geom);
create index idx_nodes_name on nodes(name);
create index idx_nodes_deleted_at on nodes(deleted_at) where deleted_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists nodes;
-- +goose StatementEnd
