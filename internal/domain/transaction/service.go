package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	accountDomain "github.com/cavejondev/finan-simples/internal/domain/account"
	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/cavejondev/finan-simples/internal/domain/database"
	"github.com/cavejondev/finan-simples/internal/domain/logger"
	subcategoryDomain "github.com/cavejondev/finan-simples/internal/domain/subcategory"
)

// Erros que podem ocorrer
var (
	ErrAccountRequired     = errors.New("account is required")
	ErrAccountNotFound     = errors.New("account not found")
	ErrSubcategoryRequired = errors.New("subcategory is required")
	ErrSubcategoryNotFound = errors.New("subcategory not found")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrAmountInvalid       = errors.New("amount must be greater than zero")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTransferAccountSame = errors.New("cannot transfer to the same account")
)

// Service representa o servico de transaction
type Service struct {
	repository         Repository
	accountService     AccountService
	subcategoryService SubcategoryService
	categoryService    CategoryService
	txManager          database.Manager
	logger             *logger.Service
}

func NewService(
	repository Repository,
	accountService AccountService,
	subcategoryService SubcategoryService,
	categoryService CategoryService,
	txManager database.Manager,
	logger *logger.Service,
) *Service {
	return &Service{
		repository:         repository,
		accountService:     accountService,
		subcategoryService: subcategoryService,
		categoryService:    categoryService,
		txManager:          txManager,
		logger:             logger,
	}
}

func (s *Service) Create(
	ctx context.Context,
	personID uuid.UUID,
	accountID uuid.UUID,
	toAccountID uuid.UUID,
	subcategoryID uuid.UUID,
	amount int64,
	description string,
	occurredAt time.Time,
) error {
	if accountID == uuid.Nil {
		return ErrAccountRequired
	}

	if amount <= 0 {
		return ErrAmountInvalid
	}

	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	account, err := s.accountService.FindByID(ctx, personID, accountID)
	if err != nil {

		if errors.Is(err, accountDomain.ErrAccountNotFound) {
			return ErrAccountNotFound
		}

		s.logger.Error(
			ctx,
			"transaction account service find by id error",
			err,
		)

		return err
	}

	if account == nil {
		return ErrAccountNotFound
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {

		s.logger.Error(
			ctx,
			"transaction begin tx error",
			err,
		)

		return err
	}

	defer tx.Rollback()

	now := time.Now()

	var desc *string
	if description != "" {
		desc = &description
	}

	// -------------------------
	// TRANSFER
	// -------------------------

	if toAccountID != uuid.Nil {

		if accountID == toAccountID {
			s.logger.Error(ctx, "transaction transfer same account", ErrTransferAccountSame)
			return ErrTransferAccountSame
		}

		toAccount, err := s.accountService.FindByID(ctx, personID, toAccountID)
		if err != nil {

			s.logger.Error(
				ctx,
				"transaction find destination account error",
				err,
			)

			return err
		}

		if toAccount == nil {
			s.logger.Error(ctx, "transaction destination account not found", ErrAccountNotFound)
			return ErrAccountNotFound
		}

		// Gera um ID único para a transferência
		transferID := uuid.New()

		expenseTransaction := &Transaction{
			PersonID:    personID,
			ID:          uuid.New(),
			AccountID:   accountID,
			TransferID:  &transferID, // link com a transferência
			Type:        category.Expense,
			Amount:      amount,
			Description: desc,
			OccurredAt:  occurredAt,
			CreatedAt:   now,
		}

		incomeTransaction := &Transaction{
			PersonID:    personID,
			ID:          uuid.New(),
			AccountID:   toAccountID,
			TransferID:  &transferID, // mesmo ID da transferência
			Type:        category.Income,
			Amount:      amount,
			Description: desc,
			OccurredAt:  occurredAt,
			CreatedAt:   now,
		}

		if err := s.repository.Create(ctx, tx, expenseTransaction); err != nil {

			s.logger.Error(
				ctx,
				"transaction repository create expense error",
				err,
			)

			return err
		}

		if err := s.repository.Create(ctx, tx, incomeTransaction); err != nil {

			s.logger.Error(
				ctx,
				"transaction repository create income error",
				err,
			)

			return err
		}

		if err := s.accountService.DecreaseBalance(ctx, tx, accountID, amount); err != nil {

			s.logger.Error(
				ctx,
				"transaction decrease balance error",
				err,
			)

			return err
		}

		if err := s.accountService.IncreaseBalance(ctx, tx, toAccountID, amount); err != nil {

			s.logger.Error(
				ctx,
				"transaction increase balance error",
				err,
			)

			return err
		}

		if err := tx.Commit(); err != nil {

			s.logger.Error(
				ctx,
				"transaction commit error",
				err,
			)

			return err
		}

		return nil
	}

	// -------------------------
	// INCOME / EXPENSE
	// -------------------------

	if subcategoryID == uuid.Nil {
		return ErrSubcategoryRequired
	}

	subcategory, err := s.subcategoryService.FindByID(ctx, personID, subcategoryID)
	if err != nil {

		if errors.Is(err, subcategoryDomain.ErrSubcategoryNotFound) {
			return ErrSubcategoryNotFound
		}

		s.logger.Error(
			ctx,
			"transaction find subcategory error",
			err,
		)

		return err
	}

	if subcategory == nil {
		return ErrSubcategoryNotFound
	}

	categoryResult, err := s.categoryService.FindByID(ctx, personID, subcategory.CategoryID)
	if err != nil {

		s.logger.Error(
			ctx,
			"transaction find category error",
			err,
		)

		return err
	}

	if categoryResult == nil {
		s.logger.Error(ctx, "transaction category not found", ErrCategoryNotFound)
		return ErrCategoryNotFound
	}

	transaction := &Transaction{
		ID:            uuid.New(),
		PersonID:      personID,
		AccountID:     accountID,
		CategoryID:    &categoryResult.ID,
		SubcategoryID: &subcategoryID,
		Type:          categoryResult.Type,
		Amount:        amount,
		Description:   desc,
		OccurredAt:    occurredAt,
		CreatedAt:     now,
	}

	if err := s.repository.Create(ctx, tx, transaction); err != nil {

		s.logger.Error(
			ctx,
			"transaction repository create error",
			err,
		)

		return err
	}

	switch categoryResult.Type {

	case category.Income:

		err = s.accountService.IncreaseBalance(ctx, tx, accountID, amount)

	case category.Expense:

		err = s.accountService.DecreaseBalance(ctx, tx, accountID, amount)
	}

	if err != nil {

		s.logger.Error(
			ctx,
			"transaction update account balance error",
			err,
		)

		return err
	}

	if err := tx.Commit(); err != nil {

		s.logger.Error(
			ctx,
			"transaction commit error",
			err,
		)

		return err
	}

	return nil
}

func (s *Service) FindByPersonID(
	ctx context.Context,
	personID uuid.UUID,
) ([]*Transaction, error) {
	transactions, err := s.repository.FindByPersonID(ctx, personID)
	if err != nil {

		s.logger.Error(
			ctx,
			"transaction repository find by person id error",
			err,
		)

		return nil, err
	}

	return transactions, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	personID uuid.UUID,
	transactionID uuid.UUID,
) (*Transaction, error) {
	if transactionID == uuid.Nil {
		s.logger.Error(ctx, "transaction invalid id", ErrTransactionNotFound)
		return nil, ErrTransactionNotFound
	}

	transaction, err := s.repository.FindByID(ctx, personID, transactionID)
	if err != nil {

		s.logger.Error(
			ctx,
			"transaction repository find by id error",
			err,
		)

		return nil, err
	}

	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}
