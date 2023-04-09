package handlers

import (
	"net/http"
	"os"

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
	user, err := db.GetAccountByEmail(&loginData.Email)
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

	// Generate JWT token
	SecretKey := os.Getenv("SECRET_KEY")
	token, err := utils.GenerateJWT(int(user.DbId), SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
