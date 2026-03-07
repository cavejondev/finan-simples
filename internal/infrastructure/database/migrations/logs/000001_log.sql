-- +goose Up
-- +goose StatementBegin

CREATE TYPE log_level AS ENUM (
    'DEBUG',
    'INFO',
    'WARN',
    'ERROR',
    'FATAL'
);

CREATE TYPE http_method AS ENUM (
    'GET',
    'POST',
    'PUT',
    'PATCH',
    'DELETE'
);

CREATE TABLE logs (
    id UUID NOT NULL,
    level log_level NOT NULL,
    message TEXT NOT NULL,

    service TEXT,
    request_id VARCHAR(100),
    user_id UUID,

    method http_method,
    path TEXT,
    status_code INTEGER,

    error TEXT,
    metadata JSONB,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT logs_pk PRIMARY KEY (id)
);

CREATE INDEX idx_logs_level ON logs(level);
CREATE INDEX idx_logs_created_at ON logs(created_at);
CREATE INDEX idx_logs_request_id ON logs(request_id);

-- +goose StatementEnd