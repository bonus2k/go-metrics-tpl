CREATE TABLE gauge
(
    id    serial PRIMARY KEY,
    name  VARCHAR(255) UNIQUE,
    value DOUBLE PRECISION NOT NULL
);