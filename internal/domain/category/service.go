package category

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
)

// Erros do serviço
var (
	ErrNameRequired     = errors.New("category name is required")
	ErrNameTooShort     = errors.New("category name too short")
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInternal = errors.New("category internal error")
	ErrNameDuplicated   = errors.New("name duplicated")
	ErrInvalidType      = errors.New("invalid category type")
)

// Erros do banco
var (
	ErrPersistenceNameDuplicated = errors.New("category name already exists")
)

// Service representa o serviço de category
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

// Create cria uma nova categoria
func (s *Service) Create(ctx context.Context, personID uuid.UUID, name string, t Type) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	// valida type
	if t != Income && t != Expense {
		return ErrInvalidType
	}

	category := &Category{
		ID:        uuid.New(),
		PersonID:  personID,
		Name:      name,
		Type:      t,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(category); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"category repository create error",
			err,
		)

		return ErrCategoryInternal
	}

	return nil
}

// Update altera nome da categoria
func (s *Service) Update(ctx context.Context, personID, categoryID uuid.UUID, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	category, err := s.repository.FindByID(personID, categoryID)
	if err != nil {

		s.logger.Error(
			ctx,
			"category repository find by id error",
			err,
		)

		return ErrCategoryInternal
	}

	if category == nil || category.ID == uuid.Nil {
		return ErrCategoryNotFound
	}

	category.Name = name

	if err := s.repository.Update(category); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"category repository update error",
			err,
		)

		return ErrCategoryInternal
	}

	return nil
}

// FindByPersonID retorna todas categorias da pessoa
func (s *Service) FindByPersonID(ctx context.Context, personID uuid.UUID) ([]*Category, error) {
	categories, err := s.repository.FindByPersonID(personID)
	if err != nil {

		s.logger.Error(
			ctx,
			"category repository find by person id error",
			err,
		)

		return nil, ErrCategoryInternal
	}

	return categories, nil
}

// FindByID retorna uma categoria específica da pessoa
func (s *Service) FindByID(ctx context.Context, personID, categoryID uuid.UUID) (*Category, error) {
	category, err := s.repository.FindByID(personID, categoryID)
	if err != nil {

		s.logger.Error(
			ctx,
			"category repository find by id error",
			err,
		)

		return nil, ErrCategoryInternal
	}

	if category == nil || category.ID == uuid.Nil {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}
