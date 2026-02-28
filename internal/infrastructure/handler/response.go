// Package handler contem as definições do pacote http
package handler

import (
	"encoding/json"
	"net/http"
)

// APIResponse é o padrão de resposta da API.
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse retorna resposta de sucesso padronizada.
func SuccessResponse(w http.ResponseWriter, status int, code string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorResponse retorna resposta de erro padronizada.
func ErrorResponse(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Code:    code,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
