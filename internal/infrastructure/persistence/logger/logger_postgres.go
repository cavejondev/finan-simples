// Package logger é o pacote de log do sistema
package logger

import (
	"github.com/jmoiron/sqlx"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
)

// Repository representa o repositório responsável por persistir logs no banco
type Repository struct {
	db *sqlx.DB
}

// NewRepository cria uma nova instância do repositório de logs
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create insere um novo log na tabela logs
func (r *Repository) Create(l *logger.Log) error {
	query := `
	INSERT INTO logs (
		id,
		level,
		message,
		service,
		request_id,
		user_id,
		method,
		path,
		error,
		metadata
	)
	VALUES (
		:id,
		:level,
		:message,
		:service,
		:request_id,
		:user_id,
		:method,
		:path,
		:error,
		:metadata
	)
	`

	_, err := r.db.NamedExec(query, l)

	return err
}
