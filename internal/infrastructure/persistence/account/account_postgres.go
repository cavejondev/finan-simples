// Package account é o pacote que faz operações de conta diretamente com banco de dados
package account

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/cavejondev/finan-simples/internal/domain/account"
	"github.com/cavejondev/finan-simples/internal/infrastructure/persistence"
)

// Repository implementa account.Repository
type Repository struct {
	db *sqlx.DB
}

// NewAccountRepository cria nova instancia do repositorio
func NewAccountRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create salva uma nova conta no banco
func (r *Repository) Create(a *account.Account) error {
	query := `
		INSERT INTO account (id, person_id, name, balance)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := r.db.QueryRowx(
		query,
		a.ID,
		a.PersonID,
		a.Name,
		a.Balance,
	).Scan(&a.CreatedAt)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == AccountPersonNameUnique {
				return account.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}

// FindByPersonID busca todas contas da pessoa
func (r *Repository) FindByPersonID(personID uuid.UUID) ([]*account.Account, error) {
	query := `
		SELECT id, person_id, name, balance, created_at, closed_at
		FROM account
		WHERE person_id = $1
	`

	var accounts []*account.Account

	err := r.db.Select(&accounts, query, personID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindByID busca conta especifica da pessoa
func (r *Repository) FindByID(personID, accountID uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, person_id, name, balance, created_at, closed_at
		FROM account
		WHERE id = $1
		AND person_id = $2
	`

	var acc account.Account

	err := r.db.Get(&acc, query, accountID, personID)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &acc, nil
}

// Update atualiza dados da conta
func (r *Repository) Update(a *account.Account) error {
	query := `
		UPDATE account
		SET name = $1
		WHERE id = $2
		AND person_id = $3
	`

	_, err := r.db.Exec(
		query,
		a.Name,
		a.ID,
		a.PersonID,
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == AccountPersonNameUnique {
				return account.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}
