-- +goose Up
-- +goose StatementBegin

CREATE TYPE transaction_type AS ENUM (
    'INCOME',
    'EXPENSE',
    'TRANSFER'
);

CREATE TABLE transaction (
    id UUID NOT NULL,
    account_id UUID NOT NULL,
    subcategory_id UUID NOT NULL,

    type transaction_type NOT NULL,
    amount BIGINT NOT NULL,

    description TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT transaction_pk PRIMARY KEY (id)
);

-- +goose StatementEnd