CREATE TABLE IF NOT EXISTS order (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    restaurant_id INT NOT NULL,
    created_at TIMESTAMP
);

CREATE INDEX order_idx_uid ON order (user_id);

CREATE TABLE IF NOT EXISTS order_products (
    order_id INT NOT NULL,
    fid INT NOT NULL,
    name VARCHAR(30) NOT NULL,
    price FLOAT NOT NULL,
    quantity INTEGER NOT NULL,

    PRIMARY KEY (order_id, fid),
    FOREIGN KEY (order_id) REFERENCES order (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_oid ON products (order_id);