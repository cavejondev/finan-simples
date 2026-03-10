// Package person é o pacote de pessoa http
package person

import "github.com/google/uuid"

// REQUESTS

// RegisterRequest representa o payload de registro.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest representa o payload de login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RESPONSES

// LoginResponse representa o retorno do login.
type LoginResponse struct {
	Token string `json:"token"`
}

// GetMeResponse representa o retorno da rota get me.
type GetMeResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
