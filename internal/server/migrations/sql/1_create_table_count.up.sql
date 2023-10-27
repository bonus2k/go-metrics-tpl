CREATE TABLE count
(
    id serial PRIMARY KEY,
    name  VARCHAR(255) UNIQUE,
    value INTEGER NOT NULL
);