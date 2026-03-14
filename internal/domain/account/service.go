package account

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/database"
	"github.com/cavejondev/finan-simples/internal/domain/logger"
)

// Erros do serviço
var (
	ErrNameRequired    = errors.New("account name is required")
	ErrNameTooShort    = errors.New("account name too short")
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountClosed   = errors.New("account is closed")
	ErrAccountInternal = errors.New("account internal error")
	ErrNameDuplicated  = errors.New("name duplicated")
)

// Erros do banco
var (
	ErrPersistenceNameDuplicated = errors.New("account name already exists")
)

// Service representa o serviço de account
type Service struct {
	repository Repository
	logger     *logger.Service
}

// NewService cria nova instancia do serviço
func NewService(
	r Repository,
	log *logger.Service,
) *Service {
	return &Service{
		repository: r,
		logger:     log,
	}
}

// Create cria uma nova conta para a pessoa
func (s *Service) Create(ctx context.Context, personID uuid.UUID, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	account := &Account{
		ID:        uuid.New(),
		PersonID:  personID,
		Name:      name,
		Balance:   0,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(account); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"account repository create error",
			err,
		)

		return ErrAccountInternal
	}

	return nil
}

// Update altera o nome da conta
func (s *Service) Update(ctx context.Context, personID, accountID uuid.UUID, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	account, err := s.repository.FindByID(personID, accountID)
	if err != nil {

		s.logger.Error(
			ctx,
			"account repository find by id error",
			err,
		)

		return ErrAccountInternal
	}

	if account == nil || account.ID == uuid.Nil {
		return ErrAccountNotFound
	}

	if account.ClosedAt != nil {
		return ErrAccountClosed
	}

	account.Name = name

	if err := s.repository.Update(account); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"account repository update error",
			err,
		)

		return ErrAccountInternal
	}

	return nil
}

// FindByPersonID retorna todas contas da pessoa
func (s *Service) FindByPersonID(ctx context.Context, personID uuid.UUID) ([]*Account, error) {
	accounts, err := s.repository.FindByPersonID(personID)
	if err != nil {

		s.logger.Error(
			ctx,
			"account repository find by person id error",
			err,
		)

		return nil, ErrAccountInternal
	}

	return accounts, nil
}

// FindByID retorna uma conta específica da pessoa
func (s *Service) FindByID(ctx context.Context, personID, accountID uuid.UUID) (*Account, error) {
	account, err := s.repository.FindByID(personID, accountID)
	if err != nil {

		s.logger.Error(
			ctx,
			"account repository find by id error",
			err,
		)

		return nil, ErrAccountInternal
	}

	if account == nil || account.ID == uuid.Nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *Service) IncreaseBalance(
	ctx context.Context,
	tx database.Tx,
	accountID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return nil
	}

	err := s.repository.IncreaseBalance(ctx, tx, accountID, amount)
	if err != nil {

		if errors.Is(err, ErrAccountNotFound) {
			return ErrAccountNotFound
		}

		s.logger.Error(
			ctx,
			"account repository increase balance error",
			err,
		)

		return ErrAccountInternal
	}

	return nil
}

func (s *Service) DecreaseBalance(
	ctx context.Context,
	tx database.Tx,
	accountID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return nil
	}

	err := s.repository.DecreaseBalance(ctx, tx, accountID, amount)
	if err != nil {

		if errors.Is(err, ErrAccountNotFound) {
			return ErrAccountNotFound
		}

		s.logger.Error(
			ctx,
			"account repository decrease balance error",
			err,
		)

		return ErrAccountInternal
	}

	return nil
}
