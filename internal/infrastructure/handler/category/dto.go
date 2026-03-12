// Package category é o pacote de category http
package category

import (
	"time"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/category"
)

//
// REQUESTS
//

// RequestCreate representa o payload para criar categoria
type RequestCreate struct {
	Name string        `json:"name"`
	Type category.Type `json:"type"`
}

// RequestUpdate representa o payload para atualizar categoria
type RequestUpdate struct {
	Name string `json:"name"`
}

//
// RESPONSES
//

// ResponseCategory representa uma categoria retornada pela API
type ResponseCategory struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Type      category.Type `json:"type"`
	CreatedAt time.Time     `json:"created_at"`
}
