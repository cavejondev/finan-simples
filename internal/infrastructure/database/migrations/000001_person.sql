-- +goose Up
-- +goose StatementBegin

CREATE SEQUENCE person_id_seq
    START 1
    INCREMENT 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE person (
    id BIGINT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT person_pk PRIMARY KEY (id),
    CONSTRAINT person_email_unique UNIQUE (email)
);

ALTER SEQUENCE person_id_seq OWNED BY person.id;

ALTER TABLE person
    ALTER COLUMN id SET DEFAULT nextval('person_id_seq');

-- +goose StatementEnd