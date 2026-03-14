package database

import (
	"context"

	domain "github.com/cavejondev/finan-simples/internal/domain/database"
	"github.com/jmoiron/sqlx"
)

type manager struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) domain.Manager {
	return &manager{
		db: db,
	}
}

func (m *manager) Begin(ctx context.Context) (domain.Tx, error) {
	sqlxTx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &tx{
		tx: sqlxTx,
	}, nil
}
