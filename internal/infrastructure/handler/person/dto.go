// Package person é o pacote de pessoa http
package person

import "github.com/google/uuid"

//
// REQUESTS
//

// RequestRegister representa o payload de registro
type RequestRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RequestLogin representa o payload de login
type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//
// RESPONSES
//

// ResponseLogin representa o retorno do login
type ResponseLogin struct {
	Token string `json:"token"`
}

// ResponsePerson representa a pessoa retornada pela API
type ResponsePerson struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// ResponseGetMe representa o retorno da rota get me.
type ResponseGetMe struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
