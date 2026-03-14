package transaction

import (
	"time"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/category"
)

//
// REQUESTS
//

type RequestCreate struct {
	AccountID     uuid.UUID `json:"account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id,omitempty"`
	SubcategoryID uuid.UUID `json:"subcategory_id"`
	Amount        int64     `json:"amount"`
	Description   string    `json:"description,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
}

//
// RESPONSES
//

type ResponseTransaction struct {
	ID            uuid.UUID     `json:"id"`
	AccountID     uuid.UUID     `json:"account_id"`
	CategoryID    uuid.UUID     `json:"category_id"`
	SubcategoryID uuid.UUID     `json:"subcategory_id"`
	Type          category.Type `json:"type"`
	Amount        int64         `json:"amount"`
	Description   *string       `json:"description,omitempty"`
	OccurredAt    time.Time     `json:"occurred_at"`
	CreatedAt     time.Time     `json:"created_at"`
}
