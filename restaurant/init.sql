CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS restaurant (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(100),
    location VARCHAR(100) NOT NULL,
    coordinate GEOGRAPHY(POINT, 4326) NOT NULL
);

CREATE INDEX idx_restaurant_location ON restaurant (location);
CREATE INDEX idx_restaurant_coord ON restaurant USING GIST(coordinate); 

CREATE TABLE IF NOT EXISTS foods (
    fid INTEGER PRIMARY KEY AUTOINCREMENT,
    rid INTEGER NOT NULL,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(100),
    price Float NOT NULL,

    FOREIGN KEY (rid) REFERENCES restaurant (id) ON DELETE CASCADE
);

CREATE INDEX idx_food_rid ON food (rid);