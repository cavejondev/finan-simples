package subcategory

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/subcategory"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler adapta requisições HTTP para o domínio subcategory.
type Handler struct {
	service *domain.Service
}

// NewHandler cria nova instância do handler.
func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

// Create cria uma subcategoria
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestCreate
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "invalid request body")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	err := h.service.Create(ctx, userID, req.CategoryID, req.Name)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeSubcategoryNameDuplicated, err.Error())

		case errors.Is(err, domain.ErrNameRequired):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameRequired, err.Error())

		case errors.Is(err, domain.ErrNameTooShort):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameTooShort, err.Error())

		case errors.Is(err, domain.ErrCategoryRequired):
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeCategoryRequired, err.Error())

		case errors.Is(err, domain.ErrCategoryNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeCategoryNotFound, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, err.Error())
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusCreated,
		CodeSubcategoryCreated,
		"Subcategory created successfully",
		nil,
	)
}

// Update altera nome da subcategoria
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestUpdate
	ctx := r.Context()

	subcategoryIDStr := chi.URLParam(r, "id")

	subcategoryID, err := uuid.Parse(subcategoryIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidSubcategoryID, "invalid subcategory id")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "invalid request body")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	err = h.service.Update(ctx, userID, subcategoryID, req.Name)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeSubcategoryNameDuplicated, err.Error())

		case errors.Is(err, domain.ErrSubcategoryNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeSubcategoryNotFound, err.Error())

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
		CodeSubcategoryUpdated,
		"Subcategory updated successfully",
		nil,
	)
}

// GetAll lista subcategorias do usuário
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	subcategories, err := h.service.FindByPersonID(ctx, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		return
	}

	res := make([]ResponseSubcategory, 0, len(subcategories))

	for _, sub := range subcategories {
		res = append(res, ResponseSubcategory{
			ID:         sub.ID,
			CategoryID: sub.CategoryID,
			Name:       sub.Name,
		})
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeSubcategoriesListed,
		"Subcategories listed",
		res,
	)
}

// Get busca subcategoria específica
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subcategoryIDStr := chi.URLParam(r, "id")

	subcategoryID, err := uuid.Parse(subcategoryIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidSubcategoryID, "invalid subcategory id")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	subcategory, err := h.service.FindByID(ctx, userID, subcategoryID)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrSubcategoryNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeSubcategoryNotFound, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeSubcategoryFound,
		"Subcategory found",
		ResponseSubcategory{
			ID:         subcategory.ID,
			CategoryID: subcategory.CategoryID,
			Name:       subcategory.Name,
		},
	)
}
