package middleware

import (
	redismethod "api/internal/redis/method"
	"context"
	"fmt"
	"net/http"
	"time"
)

type RateLimiter struct {
	RedisClient *redismethod.Redis
	RateLimit   int
	Window      time.Duration
}

func (rl *RateLimiter)Limit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		ctx := context.Background()
		key := fmt.Sprintf("rate:%s", ip)
		count, err := rl.RedisClient.R.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if count == 1 {
			rl.RedisClient.R.Expire(ctx, key, rl.Window)
		}
		if int(count) > rl.RateLimit {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
