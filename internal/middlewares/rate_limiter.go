package middlewares

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func RateLimiter(redisClient *redis.Client, limit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("rate_limit:%s", c.IP())
		ctx := context.Background()
		count, err := redisClient.Get(ctx, key).Int()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
		if count >= limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		}
		redisClient.Incr(ctx, key)
		redisClient.Expire(ctx, key, time.Second)

		return c.Next()
	}
}
