
-- +migrate Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL, -- bisa digunakan untuk login
    phone_num VARCHAR NOT NULL, -- bisa digunakan untuk login
    "avatar" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    "type" INTEGER NOT NULL,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email ON users (
    (lower(email))
);

CREATE UNIQUE INDEX users_phone_num ON users (
    (lower(phone_num))
);

CREATE TABLE user_passhashes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    passhash VARCHAR NOT NULL,
    deprecated BOOLEAN NOT NULL,
    ver INT NOT NULL, -- passhash versioning, dibutuhkan apabila ingin merubah algo passhash ketika sudah jalan.
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE register_users (
    -- id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL, -- untuk melakukan aktivasi via email
    phone_num VARCHAR NOT NULL, -- untuk melakukan aktivasi via phone (kalau tidak email)
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    token TEXT PRIMARY KEY,
    code VARCHAR(6) NOT NULL -- activation code bisa digunakan untuk aktivasi via SMS misalnya.
);

CREATE UNIQUE INDEX register_users_email ON register_users (
    (lower(email))
);

CREATE UNIQUE INDEX register_users_phone_num ON register_users (
    (lower(phone_num))
);

-- +migrate Down
DROP INDEX IF EXISTS register_users_email;
DROP INDEX IF EXISTS register_users_phone_num;
DROP INDEX IF EXISTS users_phone_num;
DROP INDEX IF EXISTS users_email;

DROP TABLE IF EXISTS user_passhashes;
DROP TABLE IF EXISTS register_users;
DROP TABLE IF EXISTS users;