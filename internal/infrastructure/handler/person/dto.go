// Package person é o pacote de pessoa http
package person

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
