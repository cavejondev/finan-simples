package subcategory

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/category"
	"github.com/cavejondev/finan-simples/internal/domain/logger"
)

// Erros do serviço
var (
	ErrNameRequired        = errors.New("subcategory name is required")
	ErrNameTooShort        = errors.New("subcategory name too short")
	ErrSubcategoryNotFound = errors.New("subcategory not found")
	ErrSubcategoryInternal = errors.New("subcategory internal error")
	ErrNameDuplicated      = errors.New("name duplicated")
	ErrCategoryRequired    = errors.New("category id is required")
	ErrCategoryNotFound    = errors.New("category not found")
)

// Erros do banco
var (
	ErrPersistenceNameDuplicated = errors.New("subcategory name already exists")
)

// Service representa o serviço de subcategory
type Service struct {
	repository      Repository
	categoryService CategoryService
	logger          *logger.Service
}

// NewService cria nova instancia do serviço
func NewService(
	r Repository,
	categoryService CategoryService,
	log *logger.Service,
) *Service {
	return &Service{
		repository:      r,
		categoryService: categoryService,
		logger:          log,
	}
}

// Create cria uma nova subcategoria para uma categoria da pessoa.
func (s *Service) Create(
	ctx context.Context,
	personID uuid.UUID,
	categoryID uuid.UUID,
	name string,
) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	if categoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	categoryEntity, err := s.categoryService.FindByID(ctx, personID, categoryID)
	if err != nil {

		// categoria não encontrada (regra de negócio)
		if errors.Is(err, category.ErrCategoryNotFound) {
			return ErrCategoryNotFound
		}

		// erro inesperado
		s.logger.Error(
			ctx,
			"category service find by id error",
			err,
		)

		return ErrSubcategoryInternal
	}

	// segurança extra: garante que realmente veio uma categoria válida
	if categoryEntity == nil || categoryEntity.ID == uuid.Nil {
		return ErrCategoryNotFound
	}

	subcategory := &Subcategory{
		ID:         uuid.New(),
		PersonID:   personID,
		CategoryID: categoryID,
		Name:       name,
	}

	if err := s.repository.Create(subcategory); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"subcategory repository create error",
			err,
		)

		return ErrSubcategoryInternal
	}

	return nil
}

// Update altera nome da subcategoria
func (s *Service) Update(
	ctx context.Context,
	personID,
	subcategoryID uuid.UUID,
	name string,
) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrNameRequired
	}

	if len(name) < 3 {
		return ErrNameTooShort
	}

	subcategory, err := s.repository.FindByID(personID, subcategoryID)
	if err != nil {

		s.logger.Error(
			ctx,
			"subcategory repository find by id error",
			err,
		)

		return ErrSubcategoryInternal
	}

	if subcategory == nil || subcategory.ID == uuid.Nil {
		return ErrSubcategoryNotFound
	}

	subcategory.Name = name

	if err := s.repository.Update(subcategory); err != nil {

		if errors.Is(err, ErrPersistenceNameDuplicated) {
			return ErrNameDuplicated
		}

		s.logger.Error(
			ctx,
			"subcategory repository update error",
			err,
		)

		return ErrSubcategoryInternal
	}

	return nil
}

// FindByPersonID retorna todas subcategorias da pessoa
func (s *Service) FindByPersonID(
	ctx context.Context,
	personID uuid.UUID,
) ([]*Subcategory, error) {
	subcategories, err := s.repository.FindByPersonID(personID)
	if err != nil {

		s.logger.Error(
			ctx,
			"subcategory repository find by person id error",
			err,
		)

		return nil, ErrSubcategoryInternal
	}

	return subcategories, nil
}

// FindByID retorna uma subcategoria específica
func (s *Service) FindByID(
	ctx context.Context,
	personID,
	subcategoryID uuid.UUID,
) (*Subcategory, error) {
	subcategory, err := s.repository.FindByID(personID, subcategoryID)
	if err != nil {

		s.logger.Error(
			ctx,
			"subcategory repository find by id error",
			err,
		)

		return nil, ErrSubcategoryInternal
	}

	if subcategory == nil || subcategory.ID == uuid.Nil {
		return nil, ErrSubcategoryNotFound
	}

	return subcategory, nil
}
