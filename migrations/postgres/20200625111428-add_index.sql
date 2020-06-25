
-- +migrate Up
CREATE INDEX pruducts_user_id ON products (user_id);
CREATE INDEX pruduct_images_product_id ON product_images (product_id);
CREATE INDEX pruduct_labels_product_id ON product_labels (product_id);
CREATE INDEX pruduct_bidders_user_id ON product_bidders (user_id);
CREATE INDEX pruduct_bidders_product_id ON product_bidders (product_id);
CREATE INDEX access_tokens_user_id ON access_tokens (user_id);
CREATE INDEX user_passhashes_user_id ON user_passhashes (user_id);
-- +migrate Down
DROP INDEX IF EXISTS pruducts_user_id;
DROP INDEX IF EXISTS pruduct_images_product_id;
DROP INDEX IF EXISTS pruduct_labels_product_id;
DROP INDEX IF EXISTS pruduct_bidders_user_id;
DROP INDEX IF EXISTS pruduct_bidders_product_id;
DROP INDEX IF EXISTS access_tokens_user_id;
DROP INDEX IF EXISTS user_passhashes_user_id;
