
-- +migrate Up
ALTER TABLE products ADD COLUMN sold BOOLEAN DEFAULT 'f';
ALTER TABLE products ADD COLUMN closed BOOLEAN DEFAULT 'f';
-- +migrate Down
ALTER TABLE products DROP COLUMN IF EXISTS sold;
ALTER TABLE products DROP COLUMN IF EXISTS closed;