package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ARTMUC/magic-video/internal/logger"
	"github.com/ARTMUC/magic-video/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

func Auth[E any](
	api huma.API,
	parseFunc func(tokenStr string, isRefresh bool) (*service.JWTClaimsWithEntity[E], error),
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			logger.Log.Error("Authorization header is empty")
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		claims, err := parseFunc(tokenString, false)
		if err != nil {
			logger.Log.Error("Failed to parse token", zap.Error(err))
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			logger.Log.Error("Failed to get expiration time", zap.Error(err))
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if expirationTime.Before(time.Now().UTC()) {
			logger.Log.Info("Token expired")
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx = huma.WithValue(ctx, "auth", claims)

		next(ctx)
	}
}
