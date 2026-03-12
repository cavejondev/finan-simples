package category

import "github.com/google/uuid"

// Repository é o repositorio de conta
type Repository interface {
	Create(category *Category) error
	Update(category *Category) error
	FindByPersonID(personID uuid.UUID) ([]*Category, error)
	FindByID(personID, categoryID uuid.UUID) (*Category, error)
}
