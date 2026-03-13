// Package subcategory é o pacote que faz operações de subcategoria diretamente com banco de dados
package subcategory

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	domain "github.com/cavejondev/finan-simples/internal/domain/subcategory"
	"github.com/cavejondev/finan-simples/internal/infrastructure/persistence"
)

// Repository implementa subcategory.Repository
type Repository struct {
	db *sqlx.DB
}

// NewSubcategoryRepository cria nova instancia do repositorio
func NewSubcategoryRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create salva nova subcategoria
func (r *Repository) Create(s *domain.Subcategory) error {
	query := `
		INSERT INTO subcategory (id, person_id, category_id, name)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		s.ID,
		s.PersonID,
		s.CategoryID,
		s.Name,
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == SubcategoryCategoryNameUnique {
				return domain.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}

// FindByPersonID busca todas subcategorias da pessoa
func (r *Repository) FindByPersonID(personID uuid.UUID) ([]*domain.Subcategory, error) {
	query := `
		SELECT id, person_id, category_id, name
		FROM subcategory
		WHERE person_id = $1
	`

	var subcategories []*domain.Subcategory

	err := r.db.Select(&subcategories, query, personID)
	if err != nil {
		return nil, err
	}

	return subcategories, nil
}

// FindByCategoryID busca subcategorias de uma categoria
func (r *Repository) FindByCategoryID(personID, categoryID uuid.UUID) ([]*domain.Subcategory, error) {
	query := `
		SELECT id, person_id, category_id, name
		FROM subcategory
		WHERE category_id = $1
		AND person_id = $2
	`

	var subcategories []*domain.Subcategory

	err := r.db.Select(&subcategories, query, categoryID, personID)
	if err != nil {
		return nil, err
	}

	return subcategories, nil
}

// FindByID busca subcategoria especifica
func (r *Repository) FindByID(personID, subcategoryID uuid.UUID) (*domain.Subcategory, error) {
	query := `
		SELECT id, person_id, category_id, name
		FROM subcategory
		WHERE id = $1
		AND person_id = $2
	`

	var sub domain.Subcategory

	err := r.db.Get(&sub, query, subcategoryID, personID)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &sub, nil
}

// Update atualiza subcategoria
func (r *Repository) Update(s *domain.Subcategory) error {
	query := `
		UPDATE subcategory
		SET name = $1
		WHERE id = $2
		AND person_id = $3
	`

	_, err := r.db.Exec(
		query,
		s.Name,
		s.ID,
		s.PersonID,
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == SubcategoryCategoryNameUnique {
				return domain.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}
