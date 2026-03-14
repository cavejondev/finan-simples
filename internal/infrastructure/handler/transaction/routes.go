package transaction

import (
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/middleware"
	"github.com/cavejondev/finan-simples/internal/infrastructure/security"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registra as rotas de transaction
func RegisterRoutes(
	r chi.Router,
	handler *Handler,
	jwtService *security.JWTService,
) {
	r.Route("/transaction", func(r chi.Router) {
		// todas rotas são privadas
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtService))

			r.Post("/", handler.Create)

			r.Get("/", handler.GetAll)

			r.Get("/{id}", handler.Get)
		})
	})
}
