// Package middleware contém middlewares HTTP
package middleware

import (
	"context"
	"net/http"
	"strings"

	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	sharedhttp "github.com/cavejondev/finan-simples/internal/infrastructure/handler"
	returncodes "github.com/cavejondev/finan-simples/internal/infrastructure/handler/returncodes"
	"github.com/cavejondev/finan-simples/internal/infrastructure/security"
)

// AuthMiddleware valida o JWT e injeta o user_id no context
func AuthMiddleware(
	jwtService *security.JWTService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, returncodes.CodeUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, returncodes.CodeUnauthorized, "invalid token format")
				return
			}

			token := parts[1]

			userID, err := jwtService.Validate(token)
			if err != nil {
				sharedhttp.ErrorResponse(w, r, http.StatusUnauthorized, returncodes.CodeUnauthorized, "user unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), contextutil.UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
