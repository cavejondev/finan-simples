// Package person é o pacote que faz operações de pessoa diretamente com banco de dados
package person

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/cavejondev/finan-simples/internal/domain/person"
	"github.com/cavejondev/finan-simples/internal/infrastructure/persistence"
)

// Repository implementa person.Repository
type Repository struct {
	db *sqlx.DB
}

// NewPersonRepository cria nova instancia do repositorio
func NewPersonRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create salva uma nova pessoa no banco
func (r *Repository) Create(p *person.Person) error {
	query := `
		INSERT INTO person (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		p.Name,
		p.Email,
		p.Password,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == persistence.PgUniqueViolation && pqErr.Constraint == EmailUnique {
				return person.ErrPersistenceEmailDuplicated
			}
		}

		return err
	}

	return nil
}

// FindByEmail busca pessoa pelo email
func (r *Repository) FindByEmail(email string) (*person.Person, error) {
	var p person.Person

	query := `
		SELECT id, name, email, password, created_at
		FROM person
		WHERE email = $1
	`

	err := r.db.Get(&p, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}
