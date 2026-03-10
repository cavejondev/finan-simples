// Package handler contem as definições do pacote http
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
)

// APIResponse é o padrão de resposta da API.
type APIResponse struct {
	RequestID string      `json:"request_id"`
	Method    string      `json:"method,omitempty"`
	Path      string      `json:"path,omitempty"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// writeResponse é o método central que monta a resposta
func writeResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code string,
	message string,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	ctx := r.Context()

	var requestID string
	if v := contextutil.GetRequestID(ctx); v != nil {
		requestID = v.String()
	}

	var method string
	if v := contextutil.GetMethod(ctx); v != nil {
		method = string(*v)
	}

	var path string
	if v := contextutil.GetPath(ctx); v != nil {
		path = *v
	}

	response := APIResponse{
		RequestID: requestID,
		Method:    method,
		Path:      path,
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	_ = json.NewEncoder(w).Encode(response)
}

// SuccessResponse retorna resposta de sucesso padronizada.
func SuccessResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code string,
	message string,
	data interface{},
) {
	writeResponse(w, r, status, code, message, data)
}

// ErrorResponse retorna resposta de erro padronizada.
func ErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code string,
	message string,
) {
	writeResponse(w, r, status, code, message, nil)
}
