-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table "country"
(
    id          serial PRIMARY KEY,
    name        varchar(50) NOT NULL,
    description varchar(200)
);

create table "company"
(
    id          serial PRIMARY KEY,
    name        varchar(50) NOT NULL,
    description varchar(200)
);

create table "country_company"
(
    company_id int NOT NULL REFERENCES "company"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    country_id int NOT NULL REFERENCES "country"(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table "user"
(
    id                serial primary key,
    name              varchar(80),
    age               smallint check (age >= 18),
    email             varchar(50) not null unique,
    password          varchar(64) not null,
    salt              varchar(12),
    registration_date timestamp,
    country           varchar(50),
    phone             varchar(50),
    surname           varchar(80)
);

create table "ticket"
(
    id      serial      PRIMARY KEY,
    user_id int         NOT NULL REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    title   varchar(50) NOT NULL,
    content text
);

create table "currency"
(
    id         serial      PRIMARY KEY,
    name       varchar(50) NOT NULL,
    code       smallint    NOT NULL    CHECK (code >= 0),
    short_name varchar(5)  NOT NULL
);

create table "portfolio"
(
    id            serial      PRIMARY KEY,
    user_id       int         NOT NULL REFERENCES "user"(id)     ON DELETE CASCADE ON UPDATE CASCADE,
    currency_id   int         NOT NULL REFERENCES "currency"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name          varchar(50) NOT NULL,
    creation_date timestamp   NOT NULL DEFAULT NOW()
);

create table "stock_category"
(
    id                       serial      PRIMARY KEY,
    name                     varchar(50) NOT NULL,
    stock_category_parent_id int         REFERENCES "stock_category"(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table "stock"
(
    id          serial      PRIMARY KEY,
    name        varchar(50) NOT NULL,
    description varchar(200),
    category_id int         NOT NULL REFERENCES "stock_category"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    currency_id int         NOT NULL REFERENCES "currency"(id)       ON DELETE CASCADE ON UPDATE CASCADE
);

create table "trading_history"
(
    id        serial    PRIMARY KEY,
    user_id   int       NOT NULL REFERENCES "user"(id)  ON DELETE CASCADE ON UPDATE CASCADE,
    stock_id  int       NOT NULL REFERENCES "stock"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    sell_time timestamp NOT NULL DEFAULT now(),
    amount    int       NOT NULL CHECK (amount > 0),
    price     int       NOT NULL CHECK (price > 0),
    type      smallint  NOT NULL DEFAULT 0
);

create table "portfolio_item"
(
    id           serial PRIMARY KEY,
    portfolio_id int    NOT NULL REFERENCES "portfolio"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stock_id     int    NOT NULL REFERENCES "stock"(id)     ON DELETE CASCADE ON UPDATE CASCADE,
    amount       int    NOT NULL CHECK (amount > 0)
);
alter table portfolio_item
    add constraint portfolio_item_pk
        unique (portfolio_id, stock_id);

create table "trading_signal"
(
    id          serial    PRIMARY KEY,
    stock_id    int       NOT NULL REFERENCES "stock"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    signal_time timestamp NOT NULL,
    price       int       NOT NULL CHECK (price > 0),
    type        smallint  NOT NULL DEFAULT 0
);

create table "stock_price"
(
    id       serial    PRIMARY KEY,
    stock_id int       NOT NULL REFERENCES "stock"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    date     timestamp NOT NULL,
    price    int       NOT NULL CHECK (price > 0)
);

create table "stock_predication"
(
    id            serial    PRIMARY KEY,
    stock_id      int       NOT NULL REFERENCES "stock"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    date          timestamp NOT NULL,
    predict_price int       NOT NULL CHECK (predict_price > 0)
);

create table "stock_category_property"
(
    id                serial       PRIMARY KEY,
    stock_category_id int          NOT NULL REFERENCES "stock_category"(id),
    name              varchar(50)  NOT NULL,
    is_required       smallint     DEFAULT 0,
    description       varchar(200)
);

create table "stock_property_value"
(
    id                         serial PRIMARY KEY,
    stock_category_property_id int    NOT NULL REFERENCES "stock_category_property"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stock_id                   int    NOT NULL REFERENCES "stock"(id)                   ON DELETE CASCADE ON UPDATE CASCADE,
    value                      text   NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
