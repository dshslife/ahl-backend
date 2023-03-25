package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/utils"
)

// Login handles the POST /auth/login endpoint
func Login(c *gin.Context) {
	// Parse request body
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user from database
	user, err := db.GetUserByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	hashpw := []byte(loginData.Password)
	pw, err := utils.HashPassword(hashpw)
	userpw := []byte(user.Password)

	// Verify password
	if !utils.VerifyPassword(pw, userpw) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// SecretKey 생성 파트는 맡겼다!

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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
	userID, err := db.VerifyToken(tokenString, secretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Set user ID in request context
	c.Set("userID", userID)

	c.Next()
}
