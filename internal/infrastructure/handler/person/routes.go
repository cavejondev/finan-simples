package person

import "github.com/go-chi/chi/v5"

// RegisterRoutes registra rotas do domínio person.
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/person", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}
