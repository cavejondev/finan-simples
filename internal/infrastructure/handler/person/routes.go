package person

import (
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/middleware"
	"github.com/cavejondev/finan-simples/internal/infrastructure/security"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registra as rotas de pessoa
func RegisterRoutes(
	r chi.Router,
	handler *Handler,
	jwtService *security.JWTService,
) {
	r.Route("/person", func(r chi.Router) {
		// publicas
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		// privadas
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtService))

			r.Get("/me", handler.Me)
		})
	})
}
