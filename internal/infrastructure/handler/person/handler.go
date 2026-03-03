package person

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/cavejondev/finan-simples/internal/domain/person"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"
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
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeInvalidBody, "Invalid request body")
		return
	}

	err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {

		switch err {

		// NAME
		case domain.ErrNameRequired:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeNameRequired, err.Error())

		case domain.ErrNameTooShort:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeNameTooShort, err.Error())

		// EMAIL
		case domain.ErrEmailRequired:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeEmailRequired, err.Error())

		case domain.ErrEmailTooShort:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeEmailTooShort, err.Error())

		case domain.ErrEmailInvalid:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeEmailInvalid, err.Error())

		// PASSWORD
		case domain.ErrPasswordRequired:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodePasswordRequired, err.Error())

		case domain.ErrPasswordTooShort:
			sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodePasswordTooShort, err.Error())

		// BUSINESS
		case domain.ErrPersonDuplicated:
			sharedhttp.ErrorResponse(w, http.StatusConflict, CodePersonDuplicated, err.Error())

		default:
			sharedhttp.ErrorResponse(w, http.StatusInternalServerError, CodeInternalError, "internal error")
		}

		return
	}

	sharedhttp.SuccessResponse(
		w,
		http.StatusCreated,
		CodeRegistred,
		"Person registered successfully",
		nil,
	)
}

// Login endpoint de autenticação.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.ErrorResponse(w, http.StatusBadRequest, CodeInvalidBody, "Invalid request body")
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			sharedhttp.ErrorResponse(w, http.StatusUnauthorized, CodeInvalidCredentials, err.Error())
		default:
			sharedhttp.ErrorResponse(w, http.StatusInternalServerError, CodeInternalError, err.Error())
		}

		return
	}

	sharedhttp.SuccessResponse(w, http.StatusOK, CodeAuthenticated, "Authenticated person", LoginResponse{Token: token})
}
