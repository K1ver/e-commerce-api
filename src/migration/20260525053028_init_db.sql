-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id uuid primary key default gen_random_uuid(),
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    username varchar(255) unique not null,
    email varchar(255) unique not null,
    password_hash text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS products (
    id uuid primary key default gen_random_uuid(),
    name varchar(255) not null,
    description text not null,
    price bigint not null check (price > 0),
    stock int not null check (stock >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS carts (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null unique references users (id) on delete cascade,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS cart_items (
    id uuid primary key default gen_random_uuid(),
    cart_id uuid not null references carts(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade,
    quantity int not null check (quantity > 0),
    unique (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    total_price bigint not null check (total_price > 0),
    status varchar(255) not null check (status in ('pending', 'paid', 'shipped', 'completed', 'canceled')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS order_items (
    id uuid primary key default gen_random_uuid(),
    order_id uuid not null references orders(id),
    product_id uuid not null references products(id) on delete cascade,
    price bigint not null check (price > 0),
    quantity int not null check (quantity > 0),
    unique(order_id, product_id)
);

CREATE TABLE IF NOT EXISTS payments (
    id uuid primary key default gen_random_uuid(),
    order_id uuid not null references orders (id),
    amount bigint not null check (amount > 0),
    status varchar(255) not null check (status in ('pending', 'success', 'failed') ),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_cart_user_id on carts(user_id);
create index idx_cart_item_cart_id on cart_items(cart_id);
create index idx_cart_item_product_id on cart_items(product_id);
create index idx_order_user_id on orders(user_id);
create index idx_order_item_order_id on order_items(order_id);
create index idx_order_item_product_id on order_items(product_id);
create index idx_payment_order_id on payments(order_id);


-- +goose Down
DROP TABLE IF EXISTS payments, order_items, orders, cart_items, carts, products, users;
