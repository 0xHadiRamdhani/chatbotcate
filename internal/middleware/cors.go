package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func CORSWithConfig(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if !allowed && len(allowedOrigins) > 0 && allowedOrigins[0] != "*" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set CORS headers
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		if len(allowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", joinStrings(allowedMethods, ", "))
		}

		if len(allowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ", "))
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string, separator string) string {
	result := ""
	for i, str := range strs {
		if i > 0 {
			result += separator
		}
		result += str
	}
	return result
}

// CORS middleware for API endpoints
func CORSForAPI() gin.HandlerFunc {
	return CORSWithConfig(
		[]string{"*"},
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		[]string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
	)
}

// CORS middleware for web endpoints
func CORSForWeb() gin.HandlerFunc {
	return CORSWithConfig(
		[]string{"http://localhost:3000", "https://yourdomain.com"},
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		[]string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
	)
}

// CORS middleware for mobile apps
func CORSForMobile() gin.HandlerFunc {
	return CORSWithConfig(
		[]string{"*"},
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		[]string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "X-API-Key"},
	)
}