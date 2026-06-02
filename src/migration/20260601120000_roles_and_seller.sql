-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role varchar(32) NOT NULL DEFAULT 'buyer'
        CHECK (role IN ('admin', 'seller', 'buyer'));

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS seller_id uuid REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products (seller_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_products_seller_id;
ALTER TABLE products DROP COLUMN IF EXISTS seller_id;
ALTER TABLE users DROP COLUMN IF EXISTS role;
