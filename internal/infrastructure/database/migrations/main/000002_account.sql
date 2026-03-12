-- +goose Up
-- +goose StatementBegin

CREATE TABLE account (
    id UUID NOT NULL,
    person_id UUID NOT NULL,
    name TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP NULL,

    CONSTRAINT account_pk PRIMARY KEY (id),

    CONSTRAINT account_person_fk
        FOREIGN KEY (person_id)
        REFERENCES person(id)
        ON DELETE RESTRICT,

    CONSTRAINT account_person_name_unique
        UNIQUE (person_id, name)
);

-- índice importante para buscar contas do usuário
CREATE INDEX account_person_idx ON account(person_id);

-- +goose StatementEnd