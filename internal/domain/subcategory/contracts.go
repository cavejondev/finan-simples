package subcategory

import (
	"context"

	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/google/uuid"
)

// Repository é o repositorio de subcategoria
type Repository interface {
	Create(subcategory *Subcategory) error
	Update(subcategory *Subcategory) error
	FindByPersonID(personID uuid.UUID) ([]*Subcategory, error)
	FindByCategoryID(personID, categoryID uuid.UUID) ([]*Subcategory, error)
	FindByID(personID, subcategoryID uuid.UUID) (*Subcategory, error)
}

type CategoryService interface {
	FindByID(ctx context.Context, personID, categoryID uuid.UUID) (*category.Category, error)
}
