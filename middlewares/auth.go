package middlewares

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/google/uuid"
	"github.com/username/schoolapp/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckAuthHeader checks if the user is authenticated
func CheckAuthHeader(c *gin.Context) {
	// /auth로 시작하는 URL 다 무시
	if strings.HasPrefix(c.Request.URL.Path, "/auth") {
		c.Next()
		return
	}
	if c.Request.URL.Path == "/publickey" {
		c.Next()
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is missing",
		})
		return
	}

	token := authHeader[7:]
	publicKey, _ := c.Get("client_key") // 이미 아래에서 확인하니까... 오류 안나겠지?
	id, err := utils.ParseJWT(&token, "user_id", publicKey.(*rsa.PublicKey))
	if err != nil {
		// 이 미친놈들이 오류 반환하다가 또 오류가 날 수도 있냐
		// c.AbortWithError()......... >:(
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	stringId := id.(string)
	userID, _ := uuid.Parse(stringId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// Add the user ID to the context for later use
	c.Set("user_id", userID)

	c.Next()
}

// VerifyToken handles middleware to verify the JWT token in the request header
func VerifyToken(c *gin.Context) {
	// /auth로 시작하는 URL 다 무시
	if strings.HasPrefix(c.Request.URL.Path, "/auth") {
		c.Next()
		return
	}
	if c.Request.URL.Path == "/publickey" {
		c.Next()
		return
	}

	// Get JWT token from request header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		c.Abort()
		return
	}

	// Verify JWT token
	publicKey, _ := c.Get("client_key")
	id, err := utils.ParseJWT(&tokenString, "user_id", publicKey.(*rsa.PublicKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	stringID := id.(string)

	userID, err := uuid.Parse(stringID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID in request context
	c.Set("user_id", userID)

	c.Next()
}

func EnforceClientJWTPubKey(c *gin.Context) {
	if c.Request.URL.Path == "/publickey" {
		c.Next()
		return
	}

	encoded := c.GetHeader("Client-Public-Key")
	if encoded == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing RSA public key"})
		c.Abort()
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if encoded == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse provided key"})
		c.Abort()
		return
	}

	block, _ := pem.Decode(decoded)

	if block.Type != "PUBLIC KEY" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Key is not RSA public key"})
		c.Abort()
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Key is not valid RSA public key"})
		c.Abort()
		return
	}

	c.Set("client_key", publicKey)

	c.Next()
}
