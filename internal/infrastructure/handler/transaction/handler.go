package transaction

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/transaction"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler adapta requisições HTTP para o domínio transaction
type Handler struct {
	service *domain.Service
}

// NewHandler cria nova instancia
func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

func getUserID(ctxValue any) (uuid.UUID, bool) {
	id, ok := ctxValue.(uuid.UUID)
	return id, ok
}

func toResponse(t *domain.Transaction) ResponseTransaction {
	return ResponseTransaction{
		ID:            t.ID,
		AccountID:     t.AccountID,
		SubcategoryID: *t.SubcategoryID,
		Type:          t.Type,
		Amount:        t.Amount,
		Description:   t.Description,
		OccurredAt:    t.OccurredAt,
		CreatedAt:     t.CreatedAt,
	}
}

//
// CREATE
//

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	var req RequestCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "invalid request body")
		return
	}

	userID, ok := getUserID(ctx.Value(contextutil.UserIDKey))
	if !ok {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user id not found")
		return
	}

	err := h.service.Create(
		ctx,
		userID,
		req.AccountID,
		req.ToAccountID,
		req.SubcategoryID,
		req.Amount,
		req.Description,
		req.OccurredAt,
	)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrAccountRequired):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusBadRequest,
				CodeAccountRequired,
				err.Error(),
			)

		case errors.Is(err, domain.ErrAccountNotFound):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusNotFound,
				CodeAccountNotFound,
				err.Error(),
			)

		case errors.Is(err, domain.ErrSubcategoryRequired):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusBadRequest,
				CodeSubcategoryRequired,
				err.Error(),
			)

		case errors.Is(err, domain.ErrSubcategoryNotFound):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusNotFound,
				CodeSubcategoryNotFound,
				err.Error(),
			)

		case errors.Is(err, domain.ErrCategoryNotFound):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusNotFound,
				CodeCategoryNotFound,
				err.Error(),
			)

		case errors.Is(err, domain.ErrAmountInvalid):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusBadRequest,
				CodeInvalidAmount,
				err.Error(),
			)

		case errors.Is(err, domain.ErrTransferAccountSame):
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusBadRequest,
				CodeTransferAccountSame,
				err.Error(),
			)

		default:
			sharedhttp.ErrorResponse(
				w, r,
				http.StatusInternalServerError,
				CodeInternalError,
				"internal error",
			)
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusCreated,
		CodeTransactionCreated,
		"Transaction created successfully",
		nil,
	)
}

// GetAll pega todas as transações
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := getUserID(ctx.Value(contextutil.UserIDKey))
	if !ok {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user id not found")
		return
	}

	transactions, err := h.service.FindByPersonID(ctx, userID)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		return
	}

	res := make([]ResponseTransaction, 0, len(transactions))

	for _, t := range transactions {
		res = append(res, toResponse(t))
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeTransactionsListed,
		"Transactions listed",
		res,
	)
}

//
// GET BY ID
//

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	transactionIDStr := chi.URLParam(r, "id")

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidTransactionID, "invalid transaction id")
		return
	}

	userID, ok := getUserID(ctx.Value(contextutil.UserIDKey))
	if !ok {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeUnauthorized, "user id not found")
		return
	}

	t, err := h.service.FindByID(ctx, userID, transactionID)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrTransactionNotFound):
			sharedhttp.ErrorResponse(w, r, http.StatusNotFound, CodeTransactionNotFound, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusOK,
		CodeTransactionFound,
		"Transaction found",
		toResponse(t),
	)
}
