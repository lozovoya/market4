CREATE EXTENSION pgcrypto;
-- товары
CREATE TABLE products
(
    id          UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    sku         TEXT NOT NULL,
    name        TEXT NOT NULL,
    uri         TEXT NOT NULL,
    description TEXT NOT NULL,
    is_active   BOOL NOT NULL,
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE prices
(
    id              BIGSERIAL PRIMARY KEY,
    sale_price      INTEGER NOT NULL,
    factory_price   INTEGER NOT NULL,
    discount_price  INTEGER NOT NULL,
    product_id      UUID REFERENCES products,
    is_active       BOOL NOT NULL,
    created         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    uri_name    TEXT UNIQUE,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE productcategory
(
    category_id  BIGINT NOT NULL REFERENCES categories,
    product_id UUID NOT NULL REFERENCES products,
    PRIMARY KEY (category_id, product_id)
);

CREATE TABLE shops
(
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    address         TEXT NOT NULL,
    lon             TEXT,
    lat             TEXT,
    working_hours   TEXT,
    created         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE productshop
(
    shop_id BIGINT NOT NULL REFERENCES shops,
    product_id UUID NOT NULL REFERENCES products,
    PRIMARY KEY (shop_id, product_id)
);

CREATE TABLE presence
(
    shop_id BIGINT NOT NULL REFERENCES shops,
    product_id UUID NOT NULL REFERENCES products,
    PRIMARY KEY (shop_id, product_id)
);

CREATE TABLE roles
(
    id BIGSERIAL PRIMARY KEY,
    c
);

CREATE TABLE users
(
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE userroles
(
    user_id BIGINT NOT NULL REFERENCES users,
    role_id BIGINT NOT NULL REFERENCES roles,
    PRIMARY KEY (user_id, role_id)
);
