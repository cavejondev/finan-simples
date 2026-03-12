package account

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/account"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler adapta requisições HTTP para o domínio account.
type Handler struct {
	service *domain.Service
}

// NewHandler cria nova instância do handler.
func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

// Create cria uma conta
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

	err := h.service.Create(ctx, userID, req.Name)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeAccountNameDuplicated, err.Error())

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
		http.StatusCreated,
		CodeAccountCreated,
		"Account created successfully",
		nil,
	)
}

// Update atualiza nome da conta
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestUpdate
	ctx := r.Context()

	accountIDStr := chi.URLParam(r, "id")

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidAccountID, "invalid account id")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "invalid request body")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	err = h.service.Update(ctx, userID, accountID, req.Name)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrNameDuplicated):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeAccountNameDuplicated, err.Error())

		case errors.Is(err, domain.ErrAccountNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeAccountNotFound, err.Error())

		case errors.Is(err, domain.ErrAccountClosed):
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodeAccountClosed, err.Error())

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
		CodeAccountUpdated,
		"Account updated successfully",
		nil,
	)
}

// GetAll lista contas do usuario
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDVal := ctx.Value(contextutil.UserIDKey)
	if userIDVal == nil {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user ID not found")
		return
	}

	userID := userIDVal.(uuid.UUID)

	accounts, err := h.service.FindByPersonID(ctx, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		return
	}

	res := make([]ResponseAccount, 0, len(accounts))

	for _, acc := range accounts {
		res = append(res, ResponseAccount{
			ID:        acc.ID,
			Name:      acc.Name,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
			ClosedAt:  acc.ClosedAt,
		})
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeAccountsListed,
		"Accounts listed",
		res,
	)
}

// Get busca conta especifica
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountIDStr := chi.URLParam(r, "id")

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidAccountID, "invalid account id")
		return
	}

	userID := ctx.Value(contextutil.UserIDKey).(uuid.UUID)

	account, err := h.service.FindByID(ctx, userID, accountID)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrAccountNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeAccountNotFound, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeAccountFound,
		"Account found",
		ResponseAccount{
			ID:        account.ID,
			Name:      account.Name,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			ClosedAt:  account.ClosedAt,
		},
	)
}
