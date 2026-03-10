-- +goose Up
-- +goose StatementBegin

CREATE TYPE category_type AS ENUM (
    'INCOME',
    'EXPENSE'
);

CREATE TABLE category (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    type category_type NOT NULL,

    CONSTRAINT category_pk PRIMARY KEY (id),
    CONSTRAINT category_name_unique UNIQUE (name)
);

-- +goose StatementEnd