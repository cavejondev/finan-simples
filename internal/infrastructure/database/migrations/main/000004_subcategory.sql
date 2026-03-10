-- +goose Up
-- +goose StatementBegin

CREATE TABLE subcategory (
    id UUID NOT NULL,
    category_id UUID NOT NULL,
    name TEXT NOT NULL,

    CONSTRAINT subcategory_pk PRIMARY KEY (id),

    CONSTRAINT subcategory_category_fk
        FOREIGN KEY (category_id)
        REFERENCES category(id)
);

-- +goose StatementEnd