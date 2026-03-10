-- +goose Up
-- +goose StatementBegin

CREATE TABLE account (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP NULL,

    CONSTRAINT account_pk PRIMARY KEY (id)
);

-- +goose StatementEnd