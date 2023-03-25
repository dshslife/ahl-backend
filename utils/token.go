package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

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
