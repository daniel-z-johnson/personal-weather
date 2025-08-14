-- +goose Up
ALTER TABLE locations ADD COLUMN temp REAL NOT NULL DEFAULT 0.0;
ALTER TABLE locations ADD COLUMN expires TEXT NOT NULL DEFAULT 0.0;

-- +goose Down
DROP TABLE locations;