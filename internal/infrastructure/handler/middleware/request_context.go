// Package middleware contém middlewares HTTP
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	"github.com/cavejondev/finan-simples/internal/infrastructure/handler/returncodes"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestMiddleware cria request_id, injeta dados no context e recupera panic
func RequestMiddleware(
	logService *logger.Service,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// gera request_id
			reqID := uuid.New()

			// injeta dados no context
			ctx := context.WithValue(r.Context(), contextutil.RequestIDKey, reqID)
			ctx = context.WithValue(ctx, contextutil.MethodKey, r.Method)
			ctx = context.WithValue(ctx, contextutil.PathKey, r.URL.Path)

			// adiciona no header
			w.Header().Set("X-Request-ID", reqID.String())

			// recover global
			defer func() {
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())

					logService.Error(
						ctx,
						"panic recovered\n"+stack,
						fmt.Errorf("%v", rec),
					)

					http.Error(
						rw,
						returncodes.CodeInternalServerError,
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
