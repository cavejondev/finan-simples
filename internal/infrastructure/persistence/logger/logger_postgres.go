package logger

import (
	"github.com/jmoiron/sqlx"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

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
		status_code,
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
		:status_code,
		:error,
		:metadata
	)
	`

	_, err := r.db.NamedExec(query, l)

	return err
}
