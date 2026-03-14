package database

import (
	"context"

	domain "github.com/cavejondev/finan-simples/internal/domain/database"
	"github.com/jmoiron/sqlx"
)

type tx struct {
	tx *sqlx.Tx
}

func (t *tx) ExecContext(ctx context.Context, q string, args ...any) (domain.Result, error) {
	return t.tx.ExecContext(ctx, q, args...)
}

func (t *tx) QueryRowContext(ctx context.Context, q string, args ...any) domain.Row {
	return t.tx.QueryRowxContext(ctx, q, args...)
}

func (t *tx) Commit() error {
	return t.tx.Commit()
}

func (t *tx) Rollback() error {
	return t.tx.Rollback()
}
