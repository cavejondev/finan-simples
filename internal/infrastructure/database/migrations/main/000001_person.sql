-- +goose Up
-- +goose StatementBegin

CREATE TABLE person (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT person_pk PRIMARY KEY (id),
    CONSTRAINT person_email_unique UNIQUE (email)
);

-- +goose StatementEnd