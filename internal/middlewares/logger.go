package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a custom logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request details
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Log request details
		log.Printf("[%s] %s %s | Status: %d | Latency: %v | IP: %s",
			method,
			path,
			c.Request.Proto,
			statusCode,
			latency,
			clientIP,
		)
	}
}
