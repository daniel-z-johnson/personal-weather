-- +goose Up
CREATE TABLE locations (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       City TEXT NOT NULL,
                       State TEXT NOT NULL DEFAULT '',
                       Country TEXT NOT NULL DEFAULT '',
                       Latitude REAL NOT NULL DEFAULT 0.0,
                       Longitude REAL NOT NULL DEFAULT 0.0
);

-- +goose Down
DROP TABLE locations;