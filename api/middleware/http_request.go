package middleware

import (
	"context"
	"net/http"
)

func ExtractRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "http_request", r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
