// Package transaction é o pacote de transação no domain
package transaction

import (
	"time"

	"github.com/google/uuid"
)

// Transaction representa a entidade de transação no dominio
type Transaction struct {
	ID            uuid.UUID `db:"id"`
	AccountID     uuid.UUID `db:"account_id"`
	SubcategoryID uuid.UUID `db:"subcategory_id"`

	Type   Type  `db:"type"`
	Amount int64 `db:"amount"`

	Description string `db:"description"`

	CreatedAt time.Time `db:"created_at"`
}
