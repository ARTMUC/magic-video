package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/ARTMUC/magic-video/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				errorReference := uuid.NewString()
				stack := debug.Stack()
				logger.Log.Error(
					"Recovered from panic",
					zap.String("request", r.RequestURI),
					zap.String("method", r.Method),
					zap.String("reference", errorReference),
					zap.Any("error", rec),
					zap.String("stack", string(stack)),
				)
				w.WriteHeader(http.StatusInternalServerError)
				body := map[string]string{
					"errorReference": errorReference,
				}
				jsonBody, _ := json.Marshal(body)
				w.Write([]byte(jsonBody))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
