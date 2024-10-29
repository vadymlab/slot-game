package middlewares

import (
	"github.com/gin-gonic/gin"
	libredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/vadymlab/slot-game/internal/config"
	"log"
)

// NewRateLimiter sets up and returns a Gin middleware for rate limiting requests.
// The rate limiter uses Redis as a store and applies limits based on the provided configuration.
//
// Parameters:
//   - config (*config.SlotConfig): Configuration structure containing rate limit settings.
//   - redisClient (*libredis.Client): Redis client instance used as the backend for the rate limiter.
//
// Returns:
//   - (gin.HandlerFunc): Gin middleware handler function to enforce rate limiting.
//
// Usage:
//
//	Add this middleware to your Gin router to restrict the rate of requests per client IP.
//	Example rate format: "5-S" (5 requests per second), "100-M" (100 requests per minute).
//
// Example:
//
//	router := gin.Default()
//	rateLimiter := NewRateLimiter(slotConfig, redisClient)
//	router.Use(rateLimiter)
func NewRateLimiter(config *config.SlotConfig, redisClient *libredis.Client) gin.HandlerFunc {

	// Parse the rate limit format from configuration (e.g., "5-S" for 5 requests per second).
	rate, err := limiter.NewRateFromFormatted(config.RateLimit)
	if err != nil {
		panic(err) // Panic on invalid rate format
	}

	// Initialize Redis-backed store with options for the limiter.
	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: "limiter", // Prefix for limiter keys in Redis
	})
	if err != nil {
		log.Fatal(err) // Log and terminate on store initialization failure
		return nil
	}

	// Create a new rate limiter with the specified rate and Redis store.
	rateLimiter := limiter.New(store, rate)

	// Return the Gin middleware handler function for rate limiting.
	return ginLimiter.NewMiddleware(rateLimiter)
}
