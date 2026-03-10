// Package category é o pacote de categoria no domain
package category

import "github.com/google/uuid"

// Category representa a entidade de categoria no dominio
type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Type Type      `db:"type"`
}
