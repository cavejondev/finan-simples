-- +goose Up
-- +goose StatementBegin

CREATE TABLE transaction (
    id UUID NOT NULL,
    person_id UUID NOT NULL,

    account_id UUID NOT NULL,
    category_id UUID,
    subcategory_id UUID,
    transfer_id UUID, -- linka transações de transferência

    type category_type NOT NULL,

    amount BIGINT NOT NULL,

    description TEXT,

    occurred_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT transaction_pk PRIMARY KEY (id),

    CONSTRAINT transaction_person_fk
        FOREIGN KEY (person_id)
        REFERENCES person(id),

    CONSTRAINT transaction_account_fk
        FOREIGN KEY (account_id)
        REFERENCES account(id),

    CONSTRAINT transaction_category_fk
        FOREIGN KEY (category_id)
        REFERENCES category(id),

    CONSTRAINT transaction_subcategory_fk
        FOREIGN KEY (subcategory_id)
        REFERENCES subcategory(id)
);

CREATE INDEX transaction_person_idx
    ON transaction(person_id);

CREATE INDEX transaction_account_idx
    ON transaction(account_id);

CREATE INDEX transaction_category_idx
    ON transaction(category_id);

CREATE INDEX transaction_subcategory_idx
    ON transaction(subcategory_id);

CREATE INDEX transaction_occurred_at_idx
    ON transaction(occurred_at);

CREATE INDEX transaction_transfer_idx
    ON transaction(transfer_id);

-- +goose StatementEnd