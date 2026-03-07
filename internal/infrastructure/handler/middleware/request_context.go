// Package middleware contém middlewares HTTP
package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"

	"github.com/cavejondev/finan-simples/internal/domain/logger"
	contextutil "github.com/cavejondev/finan-simples/internal/domain/util"
	handlererror "github.com/cavejondev/finan-simples/internal/infrastructure/handler/error"
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
			// cria wrapper para capturar status
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// gera request_id
			reqID := uuid.New().String()

			// injeta dados no context
			ctx := context.WithValue(r.Context(), contextutil.RequestIDKey, reqID)
			ctx = context.WithValue(ctx, contextutil.MethodKey, r.Method)
			ctx = context.WithValue(ctx, contextutil.PathKey, r.URL.Path)

			// adiciona no header
			rw.Header().Set("X-Request-ID", reqID)

			// recover global
			defer func() {
				if err := recover(); err != nil {

					stack := string(debug.Stack())

					logService.Error(
						ctx,
						"panic recovered",
						nil,
					)

					logService.Debug(
						ctx,
						stack,
					)

					http.Error(
						rw,
						handlererror.CodeInternalServerError,
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
