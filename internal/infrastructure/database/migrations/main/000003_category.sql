-- +goose Up
-- +goose StatementBegin

CREATE TYPE category_type AS ENUM (
    'INCOME',
    'EXPENSE'
);

CREATE TABLE category (
    id UUID NOT NULL,
    person_id UUID NOT NULL,
    name TEXT NOT NULL,
    type category_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT category_pk PRIMARY KEY (id),

    CONSTRAINT category_person_fk
        FOREIGN KEY (person_id)
        REFERENCES person(id)
        ON DELETE RESTRICT,

    CONSTRAINT category_person_name_unique
        UNIQUE (person_id, name)
);

-- índice para buscar categorias do usuário
CREATE INDEX category_person_idx
ON category(person_id);

-- +goose StatementEnd