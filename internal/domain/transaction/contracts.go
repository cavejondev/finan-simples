package transaction

import (
	"context"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/account"
	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/cavejondev/finan-simples/internal/domain/database"
	"github.com/cavejondev/finan-simples/internal/domain/subcategory"
)

// Repository representa o repositório de transaction
type Repository interface {
	Create(ctx context.Context, tx database.Tx, t *Transaction) error
	FindByPersonID(ctx context.Context, personID uuid.UUID) ([]*Transaction, error)
	FindByID(ctx context.Context, personID uuid.UUID, transactionID uuid.UUID) (*Transaction, error)
}

// AccountService define o que transaction precisa do domínio account
type AccountService interface {
	FindByID(ctx context.Context, personID, accountID uuid.UUID) (*account.Account, error)
	IncreaseBalance(ctx context.Context, tx database.Tx, accountID uuid.UUID, amount int64) error
	DecreaseBalance(ctx context.Context, tx database.Tx, accountID uuid.UUID, amount int64) error
}

// SubcategoryService define o que transaction precisa do domínio subcategory
type SubcategoryService interface {
	FindByID(ctx context.Context, personID, subcategoryID uuid.UUID) (*subcategory.Subcategory, error)
}

// CategoryService define o que transaction precisa do domínio category
type CategoryService interface {
	FindByID(ctx context.Context, personID, categoryID uuid.UUID) (*category.Category, error)
}
