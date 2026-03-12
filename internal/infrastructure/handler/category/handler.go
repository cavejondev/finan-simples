package category

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/category"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler adapta requisições HTTP para o domínio category.
type Handler struct {
	service *domain.Service
}

// NewHandler cria nova instância do handler.
func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

// Create cria uma categoria
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestCreate
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "Invalid request body")
		return
	}

	userIDVal := ctx.Value(contextutil.UserIDKey)
	if userIDVal == nil {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user ID not found")
		return
	}

	userID := userIDVal.(uuid.UUID)

	err := h.service.Create(ctx, userID, req.Name, req.Type)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeCategoryNameDuplicated, err.Error())

		case errors.Is(err, domain.ErrNameRequired):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameRequired, err.Error())

		case errors.Is(err, domain.ErrNameTooShort):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameTooShort, err.Error())

		case errors.Is(err, domain.ErrInvalidType):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidCategoryType, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusCreated,
		CodeCategoryCreated,
		"Category created successfully",
		nil,
	)
}

// Update atualiza nome da categoria
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestUpdate
	ctx := r.Context()

	categoryIDStr := chi.URLParam(r, "id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidCategoryID, "invalid category id")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "invalid request body")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	err = h.service.Update(ctx, userID, categoryID, req.Name)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeCategoryNameDuplicated, err.Error())

		case errors.Is(err, domain.ErrCategoryNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeCategoryNotFound, err.Error())

		case errors.Is(err, domain.ErrNameRequired):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameRequired, err.Error())

		case errors.Is(err, domain.ErrNameTooShort):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameTooShort, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeCategoryUpdated,
		"Category updated successfully",
		nil,
	)
}

// GetAll lista categorias do usuario
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDVal := ctx.Value(contextutil.UserIDKey)
	if userIDVal == nil {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user ID not found")
		return
	}

	userID := userIDVal.(uuid.UUID)

	categories, err := h.service.FindByPersonID(ctx, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		return
	}

	res := make([]ResponseCategory, 0, len(categories))

	for _, cat := range categories {
		res = append(res, ResponseCategory{
			ID:        cat.ID,
			Name:      cat.Name,
			Type:      cat.Type,
			CreatedAt: cat.CreatedAt,
		})
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeCategoriesListed,
		"Categories listed",
		res,
	)
}

// Get busca categoria especifica
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categoryIDStr := chi.URLParam(r, "id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidCategoryID, "invalid category id")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	category, err := h.service.FindByID(ctx, userID, categoryID)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrCategoryNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeCategoryNotFound, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeCategoryFound,
		"Category found",
		ResponseCategory{
			ID:        category.ID,
			Name:      category.Name,
			Type:      category.Type,
			CreatedAt: category.CreatedAt,
		},
	)
}
