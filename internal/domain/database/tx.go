package database

import "context"

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) Row

	Commit() error
	Rollback() error
}

type Row interface {
	Scan(dest ...any) error
}

type Result interface {
	RowsAffected() (int64, error)
}

type Manager interface {
	Begin(ctx context.Context) (Tx, error)
}
