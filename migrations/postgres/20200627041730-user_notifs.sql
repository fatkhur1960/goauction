
-- +migrate Up
CREATE TABLE user_notifs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title VARCHAR(60) NOT NULL,
    content TEXT NOT NULL,
    notif_type INT NOT NULL,
    "target" INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);
CREATE INDEX user_notifs_user_id ON user_notifs (user_id);
-- +migrate Down
DROP INDEX IF EXISTS user_notifs_user_id;
DROP TABLE IF EXISTS user_notifs;