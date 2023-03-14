CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    password VARCHAR(255) NOT NULL,
    UNIQUE(name)
);

CREATE INDEX idx_name ON users (name);