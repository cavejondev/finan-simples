package account

import "github.com/google/uuid"

// Repository é o repositorio de conta
type Repository interface {
	Create(account *Account) error
	FindByPersonID(personID uuid.UUID) ([]*Account, error)
	FindByID(personID, accountID uuid.UUID) (*Account, error)
	Update(account *Account) error
}
