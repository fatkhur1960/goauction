
-- +migrate Up
ALTER TABLE user_notifs ADD COLUMN "read" BOOLEAN NOT NULL DEFAULT 'f';
-- +migrate Down
ALTER TABLE user_notifs DROP COLUMN IF EXISTS "read";
