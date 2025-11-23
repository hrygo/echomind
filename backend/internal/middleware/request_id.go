package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextRequestIDKey is the key for storing request ID in Gin context
	ContextRequestIDKey = "request_id"
	// HeaderRequestID is the header name for request ID
	HeaderRequestID = "X-Request-ID"
)

// RequestID is a middleware that generates or extracts a request ID for each request
// and adds it to the response headers and Gin context for logging purposes
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is provided in request header
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			// Generate new request ID if not provided
			requestID = uuid.New().String()
		}

		// Store in context for logging
		c.Set(ContextRequestIDKey, requestID)

		// Add to response headers
		c.Header(HeaderRequestID, requestID)

		c.Next()
	}
}

// GetRequestIDFromContext retrieves the Request ID from the Gin context
func GetRequestIDFromContext(c *gin.Context) (string, bool) {
	requestID, ok := c.Get(ContextRequestIDKey)
	if !ok {
		return "", false
	}

	requestIDStr, ok := requestID.(string)
	return requestIDStr, ok
}
