package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RateLimiter struct {
	redis      *redis.Client
	rateLimit  int
	windowSize time.Duration
}

func NewRateLimiter(redis *redis.Client, rateLimit int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:      redis,
		rateLimit:  rateLimit,
		windowSize: windowSize,
	}
}

func (rl *RateLimiter) IsAllowed(key string) (bool, error) {
	ctx := rl.redis.Context()
	now := time.Now()
	windowStart := now.Add(-rl.windowSize)

	// Use a sliding window approach
	count, err := rl.redis.ZCount(ctx, key, windowStart.Unix(), now.Unix()).Result()
	if err != nil {
		return false, err
	}

	if count >= int64(rl.rateLimit) {
		return false, nil
	}

	// Add current request to the window
	_, err = rl.redis.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now.Unix()),
		Member: now.Unix(),
	}).Result()
	if err != nil {
		return false, err
	}

	// Clean up old entries
	_, err = rl.redis.ZRemRangeByScore(ctx, key, "-inf", windowStart.Unix()).Result()
	if err != nil {
		return false, err
	}

	// Set expiration on the key
	_, err = rl.redis.Expire(ctx, key, rl.windowSize).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID)
		clientID := c.ClientIP()
		
		// If user is authenticated, use user ID instead
		if userID, exists := c.Get("user_id"); exists {
			clientID = userID.(uuid.UUID).String()
		}

		// Create rate limiter instance (this would typically be injected)
		redisClient := c.MustGet("redis").(*redis.Client)
		rateLimiter := NewRateLimiter(redisClient, 60, time.Minute)

		// Check rate limit
		allowed, err := rateLimiter.IsAllowed("rate_limit:" + clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": time.Minute,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimiterByIP(rateLimit int, windowSize time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// Get Redis from context
		redisClient := c.MustGet("redis").(*redis.Client)
		rateLimiter := NewRateLimiter(redisClient, rateLimit, windowSize)

		allowed, err := rateLimiter.IsAllowed("rate_limit_ip:" + clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": windowSize,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimiterByUser(rateLimit int, windowSize time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get Redis from context
		redisClient := c.MustGet("redis").(*redis.Client)
		rateLimiter := NewRateLimiter(redisClient, rateLimit, windowSize)

		allowed, err := rateLimiter.IsAllowed("rate_limit_user:" + userID.(uuid.UUID).String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": windowSize,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Rate limiter for specific endpoints
func RateLimiterForEndpoint(endpoint string, rateLimit int, windowSize time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.ClientIP()
		
		// If user is authenticated, use user ID
		if userID, exists := c.Get("user_id"); exists {
			clientID = userID.(uuid.UUID).String()
		}

		// Get Redis from context
		redisClient := c.MustGet("redis").(*redis.Client)
		rateLimiter := NewRateLimiter(redisClient, rateLimit, windowSize)

		allowed, err := rateLimiter.IsAllowed("rate_limit_endpoint:" + endpoint + ":" + clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded for this endpoint",
				"retry_after": windowSize,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}