package account

import (
	"context"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/database"
)

// Repository é o repositorio de conta
type Repository interface {
	Create(account *Account) error

	FindByPersonID(personID uuid.UUID) ([]*Account, error)

	FindByID(personID, accountID uuid.UUID) (*Account, error)

	Update(account *Account) error

	IncreaseBalance(
		ctx context.Context,
		tx database.Tx,
		accountID uuid.UUID,
		amount int64,
	) error

	DecreaseBalance(
		ctx context.Context,
		tx database.Tx,
		accountID uuid.UUID,
		amount int64,
	) error
}
