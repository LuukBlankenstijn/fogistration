-- Write your migrate up statements here
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    ip VARCHAR NOT NULL UNIQUE,
    last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

---- create above / drop below ----

DROP TABLE IF EXISTS clients;
