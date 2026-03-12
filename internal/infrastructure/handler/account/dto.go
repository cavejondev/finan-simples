// Package account é o pacote de account http
package account

import (
	"time"

	"github.com/google/uuid"
)

//
// REQUESTS
//

// RequestCreate representa o payload para criar conta
type RequestCreate struct {
	Name string `json:"name"`
}

// RequestUpdate representa o payload para atualizar conta
type RequestUpdate struct {
	Name string `json:"name"`
}

//
// RESPONSES
//

// ResponseAccount representa uma conta retornada pela API
type ResponseAccount struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Balance   int64      `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}
