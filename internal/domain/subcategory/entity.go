// Package subcategory é o pacote de subcategoria no domain
package subcategory

import "github.com/google/uuid"

// Subcategory representa a entidade de subcategoria no dominio
type Subcategory struct {
	ID         uuid.UUID `db:"id"`
	CategoryID uuid.UUID `db:"category_id"`
	Name       string    `db:"name"`
}
