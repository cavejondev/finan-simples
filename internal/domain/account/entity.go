// Package account é o pacote de conta no domain
package account

import (
	"time"

	"github.com/google/uuid"
)

// Account representa a entidade de conta no dominio
type Account struct {
	ID        uuid.UUID  `db:"id"`
	PersonID  uuid.UUID  `db:"person_id"`
	Name      string     `db:"name"`
	Balance   int64      `db:"balance"`
	CreatedAt time.Time  `db:"created_at"`
	ClosedAt  *time.Time `db:"closed_at"`
}
