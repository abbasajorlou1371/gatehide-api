package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (in production)
		// c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Control referrer information
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}
