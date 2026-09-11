-- +goose Up
CREATE TABLE persons (
    id      SERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    age     INTEGER,
    address TEXT,
    work    TEXT
);

-- +goose Down
DROP TABLE persons;
