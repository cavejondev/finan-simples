// Package transaction é o pacote de transação no domain
package transaction

import (
	"time"

	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/google/uuid"
)

// Transaction representa a entidade de transação no domínio.
type Transaction struct {
	ID uuid.UUID `db:"id"`

	PersonID uuid.UUID `db:"person_id"`

	AccountID     uuid.UUID  `db:"account_id"`
	CategoryID    *uuid.UUID `db:"category_id"`    // pode ser nulo
	SubcategoryID *uuid.UUID `db:"subcategory_id"` // pode ser nulo
	TransferID    *uuid.UUID `db:"transfer_id"`    // link para transferências

	Type   category.Type `db:"type"`
	Amount int64         `db:"amount"`

	Description *string `db:"description"`

	OccurredAt time.Time  `db:"occurred_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
