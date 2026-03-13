// Package subcategory é o pacote de subcategory http
package subcategory

import "github.com/google/uuid"

//
// REQUESTS
//

// RequestCreate representa o payload para criar subcategoria
type RequestCreate struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
}

// RequestUpdate representa o payload para atualizar subcategoria
type RequestUpdate struct {
	Name string `json:"name"`
}

//
// RESPONSES
//

// ResponseSubcategory representa uma subcategoria retornada pela API
type ResponseSubcategory struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
}
