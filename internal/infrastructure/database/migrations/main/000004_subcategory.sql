-- +goose Up
-- +goose StatementBegin

CREATE TABLE subcategory (
    id UUID NOT NULL,
    person_id UUID NOT NULL,
    category_id UUID NOT NULL,
    name TEXT NOT NULL,

    CONSTRAINT subcategory_pk PRIMARY KEY (id),

    CONSTRAINT subcategory_person_fk
        FOREIGN KEY (person_id)
        REFERENCES person(id)
        ON DELETE RESTRICT,

    CONSTRAINT subcategory_category_fk
        FOREIGN KEY (category_id)
        REFERENCES category(id)
        ON DELETE RESTRICT,

    CONSTRAINT subcategory_category_name_unique
        UNIQUE (category_id, name)
);

CREATE INDEX subcategory_person_idx
ON subcategory(person_id);

CREATE INDEX subcategory_category_idx
ON subcategory(category_id);

-- +goose StatementEnd