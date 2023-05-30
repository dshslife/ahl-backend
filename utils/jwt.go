package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func EncryptJWT(toEncrypt interface{}, claim string, secretKey string) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create a new claims instance
	claims := jwt.MapClaims{}
	claims[claim] = toEncrypt
	claims["exp"] = expirationTime

	// Create a new token instance using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// DecryptJWT verifies a JWT token and decrypts its contents
func DecryptJWT(encrypted *string, claim string) (string, error) {
	// Verify the JWT encrypted
	SecretKey := os.Getenv("SECRET_KEY")
	VerifiedToken, err := VerifyJWT(*encrypted, SecretKey)
	if err != nil {
		return "", err
	}

	claims := VerifiedToken.Claims.(jwt.MapClaims)

	// Extract the user ID from the encrypted claims
	userID, ok := claims[claim].(string)
	if !ok {
		return "", fmt.Errorf("error: extracting user ID from encrypted")
	}

	return userID, nil
}

func VerifyJWT(tokenString string, secretKey string) (*jwt.Token, error) {
	// Define the expected signing method and secret key
	signingMethod := jwt.SigningMethodHS256
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}

	// Parse the JWT token string
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Verify the token signature and expiration time
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if token.Method != signingMethod {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	exp, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration time")
	}
	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Token is valid, return it
	return token, nil
}
