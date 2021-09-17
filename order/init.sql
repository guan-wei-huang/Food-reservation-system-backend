CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    created_at TIMESTAMP
);

CREATE INDEX order_idx_uid ON orders (user_id);

CREATE TABLE IF NOT EXISTS order_products (
    order_id SERIAL NOT NULL,
    fid INT NOT NULL,
    name VARCHAR(30) NOT NULL,
    price FLOAT NOT NULL,
    quantity INTEGER NOT NULL,

    PRIMARY KEY (order_id, fid),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_oid ON order_products (order_id);