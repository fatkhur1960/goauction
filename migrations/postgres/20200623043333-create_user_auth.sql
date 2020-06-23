
-- +migrate Up
CREATE TABLE access_tokens (
    token TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created TIMESTAMP NOT NULL,
    valid_thru TIMESTAMP NOT NULL
);

CREATE INDEX idx_access_tokens_user_id ON access_tokens (
    (user_id)
);
-- +migrate Down
DROP INDEX IF EXISTS  idx_access_tokens_account_id;
DROP TABLE IF EXISTS access_tokens;