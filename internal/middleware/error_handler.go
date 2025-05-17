package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response format
type ErrorResponse struct {
	Code    string      `json:"code"`              // Machine-readable error code
	Message string      `json:"message"`           // Human-readable error message
	Details interface{} `json:"details,omitempty"` // Optional additional details
}

// ErrorHandler handles errors in a consistent way
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error processing request: %v", err)

			// Check if a response has already been written
			if c.Writer.Written() {
				return
			}

			// Return a standardized error response to the client
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "internal_server_error",
				Message: "The server encountered an unexpected error while processing your request",
				Details: err.Error(),
			})
		}
	}
}