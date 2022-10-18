package middleware

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"money_share/pkg/controller"
	"money_share/pkg/database"
	"net/http"
	"time"
)

var ctx = context.Background()

func RateLimit(rateLimit int64, duration int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			key := "RATE_LIMIT_COUNT_" + ip
			err := increaseRequestCount(key, rateLimit, duration)
			if err != nil {
				controller.ResponseError(w, "Too many requests, rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func increaseRequestCount(key string, rateLimit int64, duration int64) error {
	err := database.Redis.DB.Watch(ctx, func(tx *redis.Tx) error {
		tx.SetNX(ctx, key, 0, time.Duration(duration)*time.Second)
		count, err := tx.Incr(ctx, key).Result()
		if count > rateLimit {
			err = errors.New("rate limit exceeded")
		}
		return err
	}, key)
	return err
}
