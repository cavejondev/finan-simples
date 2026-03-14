// Package transaction é o pacote que faz operações de transaction diretamente com banco de dados
package transaction

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/cavejondev/finan-simples/internal/domain/database"
	domain "github.com/cavejondev/finan-simples/internal/domain/transaction"
)

// Repository implementa transaction.Repository
type Repository struct {
	db *sqlx.DB
}

// NewTransactionRepository cria nova instancia
func NewTransactionRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create cria nova transação (sempre dentro de tx)
func (r *Repository) Create(
	ctx context.Context,
	tx database.Tx,
	t *domain.Transaction,
) error {
	query := `
		INSERT INTO transaction (
			id,
			person_id,
			account_id,
			category_id,
			subcategory_id,
			transfer_id,
			type,
			amount,
			description,
			occurred_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING created_at
	`

	return tx.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.PersonID,
		t.AccountID,
		t.CategoryID,
		t.SubcategoryID,
		t.TransferID,
		t.Type,
		t.Amount,
		t.Description,
		t.OccurredAt,
	).Scan(&t.CreatedAt)
}

// FindByPersonID lista todas transações da pessoa
func (r *Repository) FindByPersonID(
	ctx context.Context,
	personID uuid.UUID,
) ([]*domain.Transaction, error) {
	query := `
		SELECT
			id,
			person_id,
			account_id,
			category_id,
			subcategory_id,
			transfer_id,
			type,
			amount,
			description,
			occurred_at,
			created_at,
			updated_at
		FROM transaction
		WHERE person_id = $1
		ORDER BY occurred_at DESC
	`

	var transactions []*domain.Transaction

	err := r.db.SelectContext(
		ctx,
		&transactions,
		query,
		personID,
	)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindByID busca uma transação específica
func (r *Repository) FindByID(
	ctx context.Context,
	personID uuid.UUID,
	transactionID uuid.UUID,
) (*domain.Transaction, error) {
	query := `
		SELECT
			id,
			person_id,
			account_id,
			category_id,
			subcategory_id,
			transfer_id,
			type,
			amount,
			description,
			occurred_at,
			created_at,
			updated_at
		FROM transaction
		WHERE id = $1
		AND person_id = $2
	`

	var t domain.Transaction

	err := r.db.GetContext(
		ctx,
		&t,
		query,
		transactionID,
		personID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
