package middleware

import (
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	limiterMiddleware "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(next http.Handler) http.Handler {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return limiterMiddleware.NewMiddleware(instance).Handler(next)
}
