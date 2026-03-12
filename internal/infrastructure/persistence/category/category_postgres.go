// Package category é o pacote que faz operações de categoria diretamente com banco de dados
package category

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/cavejondev/finan-simples/internal/infrastructure/persistence"
)

// Repository implementa category.Repository
type Repository struct {
	db *sqlx.DB
}

// NewCategoryRepository cria nova instancia do repositorio
func NewCategoryRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create salva uma nova categoria no banco
func (r *Repository) Create(c *category.Category) error {
	query := `
		INSERT INTO category (id, person_id, name, type)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := r.db.QueryRowx(
		query,
		c.ID,
		c.PersonID,
		c.Name,
		c.Type,
	).Scan(&c.CreatedAt)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == CategoryPersonNameUnique {
				return category.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}

// FindByPersonID busca todas categorias da pessoa
func (r *Repository) FindByPersonID(personID uuid.UUID) ([]*category.Category, error) {
	query := `
		SELECT id, person_id, name, type, created_at
		FROM category
		WHERE person_id = $1
	`

	var categories []*category.Category

	err := r.db.Select(&categories, query, personID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// FindByID busca categoria especifica da pessoa
func (r *Repository) FindByID(personID, categoryID uuid.UUID) (*category.Category, error) {
	query := `
		SELECT id, person_id, name, type, created_at
		FROM category
		WHERE id = $1
		AND person_id = $2
	`

	var cat category.Category

	err := r.db.Get(&cat, query, categoryID, personID)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &cat, nil
}

// Update atualiza dados da categoria
func (r *Repository) Update(c *category.Category) error {
	query := `
		UPDATE category
		SET name = $1
		WHERE id = $2
		AND person_id = $3
	`

	_, err := r.db.Exec(
		query,
		c.Name,
		c.ID,
		c.PersonID,
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation &&
				pqErr.Constraint == CategoryPersonNameUnique {
				return category.ErrPersistenceNameDuplicated
			}
		}

		return err
	}

	return nil
}
