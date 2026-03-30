-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS postgis;

ALTER TABLE nodes ADD COLUMN geom_geog geography(Point, 4326);

UPDATE nodes
SET geom_geog = ST_SetSRID(ST_MakePoint(geom[0], geom[1]), 4326)::geography
WHERE geom IS NOT NULL;

ALTER TABLE nodes DROP COLUMN geom;
ALTER TABLE nodes RENAME COLUMN geom_geog TO geom;

DROP INDEX IF EXISTS idx_nodes_geom;
CREATE INDEX idx_nodes_geom ON nodes USING GIST (geom);

ALTER TABLE nodes ALTER COLUMN geom SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nodes ADD COLUMN geom_point point;

UPDATE nodes
SET geom_point = point(ST_X(geom::geometry), ST_Y(geom::geometry))
WHERE geom IS NOT NULL;

ALTER TABLE nodes DROP COLUMN geom;
ALTER TABLE nodes RENAME COLUMN geom_point TO geom;

DROP INDEX IF EXISTS idx_nodes_geom;
CREATE INDEX idx_nodes_geom ON nodes USING spgist(geom);
ALTER TABLE nodes ALTER COLUMN geom SET NOT NULL;
-- +goose StatementEnd
