package person

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/person"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/returncodes"
	"github.com/google/uuid"
)

// Handler adapta requisições HTTP para o domínio person.
type Handler struct {
	service *domain.Service
}

// NewHandler cria nova instância do handler.
func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

// Register endpoint de criação de usuário.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestRegister
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "Invalid request body")
		return
	}

	err := h.service.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {

		switch err {

		// NAME
		case domain.ErrNameRequired:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameRequired, err.Error())

		case domain.ErrNameTooShort:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeNameTooShort, err.Error())

		// EMAIL
		case domain.ErrEmailRequired:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeEmailRequired, err.Error())

		case domain.ErrEmailTooShort:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeEmailTooShort, err.Error())

		case domain.ErrEmailInvalid:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeEmailInvalid, err.Error())

		// PASSWORD
		case domain.ErrPasswordRequired:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodePasswordRequired, err.Error())

		case domain.ErrPasswordTooShort:
			sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodePasswordTooShort, err.Error())

		// BUSINESS
		case domain.ErrPersonDuplicated:
			sharedhttp.ErrorResponse(w, r, http.StatusConflict, CodePersonDuplicated, err.Error())

		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		r,
		http.StatusCreated,
		CodeRegistred,
		"Person registered successfully",
		nil,
	)
}

// Login endpoint de autenticação.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RequestLogin
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, r, http.StatusBadRequest, CodeInvalidBody, "Invalid request body")
		return
	}

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, CodeInvalidCredentials, err.Error())
		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, CodeInternalError, err.Error())
		}

		return
	}

	sharedhttp.SuccessResponse(w, r, http.StatusOK, CodeAuthenticated, "Authenticated person", ResponseLogin{Token: token})
}

// Me pega os dados do usuario logado
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDVal := ctx.Value(contextutil.UserIDKey)

	if userIDVal == nil {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, returncodes.CodeUnauthorized, "user ID not found in context")
		return
	}

	// já pega como uuid.UUID direto
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, returncodes.CodeUnauthorized, "invalid user ID type")
		return
	}

	// busca a pessoa pelo ID
	person, err := h.service.FindByID(ctx, userID)
	if err != nil {
		switch {
		default:
			sharedhttp.ErrorResponse(w, r, http.StatusInternalServerError, returncodes.CodeUnauthorized, "person internal server error")
		}
		return
	}

	sharedhttp.SuccessResponse(w, r, http.StatusOK, CodeAuthenticated, "Authenticated person", ResponseGetMe{
		ID:    person.ID,
		Name:  person.Name,
		Email: person.Email,
	})
}
