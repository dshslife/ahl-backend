package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
)

// Auth middleware checks if the user is authenticated
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}

		token := authHeader[7:]
		userID, err := db.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Add the user ID to the context for later use
		c.Set("user_id", userID)

		c.Next()
	}
}

// VerifyToken handles middleware to verify the JWT token in the request header
func VerifyToken(c *gin.Context) {
	// Get JWT token from request header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		c.Abort()
		return
	}

	// Verify JWT token
	userID, err := db.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Set user ID in request context
	c.Set("userID", userID)

	c.Next()
}
