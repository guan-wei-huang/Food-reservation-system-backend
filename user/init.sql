CREATE TABLE IF NOT EXISTS user (
    id SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    password VARCHAR(20) NOT NULL,
    UNIQUE(name)
);

CREATE INDEX idx_name ON user (name);