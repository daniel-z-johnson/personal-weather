-- +goose Up
CREATE TABLE locations (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       City TEXT NOT NULL,
                       State TEXT NOT NULL DEFAULT '',
                       Country TEXT NOT NULL DEFAULT '',
                       Latitude TEXT NOT NULL DEFAULT '',
                       Longitude TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE locations;