CREATE SCHEMA gomarket;

CREATE TABLE gomarket.users (
    id              BIGSERIAL   PRIMARY KEY,
    role            VARCHAR(20) NOT NULL CHECK (role IN ('buyer', 'seller')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ,
    name            VARCHAR(100) NOT NULL CHECK(char_length(name) BETWEEN 3 and 100),
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'deleted')),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    phone           VARCHAR(15) UNIQUE CHECK (
        phone ~ '^\+[0-9]+$'
        AND 
        char_length(phone) BETWEEN 10 AND 15
    )
);

CREATE TABLE gomarket.products (
    id          BIGSERIAL   PRIMARY KEY,
    price       NUMERIC(12, 2) NOT NULL CHECK(price >= 0),
    category    VARCHAR(50) NOT NULL, 
    name        VARCHAR(100) NOT NULL CHECK(char_length(name) BETWEEN 3 and 100),
    status      VARCHAR(50) NOT NULL CHECK (status IN ('active', 'deleted')),
    updated_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    stock       INTEGER NOT NULL CHECK(stock >= 0),
    seller_id   BIGINT NOT NULL REFERENCES gomarket.users(id)
);

CREATE TABLE gomarket.cart_items (
    id          BIGSERIAL   PRIMARY KEY,
    quantity    INTEGER NOT NULL CHECK(quantity > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ,
    buyer_id    BIGINT NOT NULL REFERENCES gomarket.users(id), 
    product_id  BIGINT NOT NULL REFERENCES gomarket.products(id),

    UNIQUE (buyer_id, product_id)
);

CREATE TABLE gomarket.orders (
    id          BIGSERIAL   PRIMARY KEY, 
    buyer_id    BIGINT NOT NULL REFERENCES gomarket.users(id),
    status VARCHAR(50) NOT NULL
        CHECK (
            status IN (
                'created',
                'assembling',
                'to_sorting_center',
                'to_pickup_point',
                'delivered',
                'cancelled'
            )
        ),
    price       NUMERIC(12, 2) NOT NULL CHECK(price >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);

CREATE TABLE gomarket.order_items (
    id          BIGSERIAL   PRIMARY KEY,
    seller_id   BIGINT NOT NULL REFERENCES gomarket.users(id),
    order_id   BIGINT NOT NULL REFERENCES gomarket.orders(id),
    product_id  BIGINT NOT NULL REFERENCES gomarket.products(id),
    status VARCHAR(50) NOT NULL
        CHECK (
            status IN (
                'created',
                'delivered',
                'rejected',
                'cancelled'
            )
        ),
    quantity    INTEGER NOT NULL CHECK(quantity > 0),
    price       NUMERIC(12, 2) NOT NULL CHECK(price >= 0),
    category    VARCHAR(50) NOT NULL, 
    name        VARCHAR(100) NOT NULL CHECK(char_length(name) BETWEEN 3 and 100),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);