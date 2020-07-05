
-- +migrate Up
CREATE TABLE IF NOT EXISTS user_connect (
  user_id  BIGINT  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider_name VARCHAR(50) NOT NULL, -- eg: android, apple
  app_id  VARCHAR(500)  NOT NULL,
  PRIMARY KEY (user_id)
);
-- +migrate Down
DROP TABLE IF EXISTS user_connect;