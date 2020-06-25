
-- +migrate Up
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_name VARCHAR(100) NOT NULL,
    "desc" TEXT NOT NULL,
    condition INT NOT NULL DEFAULT 0, -- 0: SECOND / 1: NEW
    condition_avg DOUBLE PRECISION NOT NULL, -- condition average in % max 100
    start_price DOUBLE PRECISION NOT NULL,
    bid_multpl DOUBLE PRECISION NOT NULL,
    closed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL
);

CREATE TABLE product_labels (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    "name" TEXT NOT NULL
);

CREATE TABLE product_bidders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    bid_price DOUBLE PRECISION NOT NULL,
    winner BOOLEAN NOT NULL DEFAULT 'f',
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS product_bidders;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_labels;
DROP TABLE IF EXISTS products;