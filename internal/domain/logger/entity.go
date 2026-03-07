package logger

import (
	"time"

	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	"github.com/google/uuid"
)

// Log é a estrutura de log
type Log struct {
	ID        uuid.UUID           `db:"id"`
	Level     Level               `db:"level"`
	Message   string              `db:"message"`
	Service   *string             `db:"service"`
	RequestID *string             `db:"request_id"`
	UserID    *uuid.UUID          `db:"user_id"`
	Method    *contextutil.Method `db:"method"`
	Path      *string             `db:"path"`
	Error     *string             `db:"error"`
	Metadata  []byte              `db:"metadata"`
	CreatedAt time.Time           `db:"created_at"`
}
