package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/utils"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

var Oauth2Application *oauth2.Config = nil

func OAuthRedirect(ctx *gin.Context) {
	url := Oauth2Application.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func OnOAuth(ctx *gin.Context) {
	code := ctx.Query("code")

	tok, err := Oauth2Application.Exchange(context.TODO(), code)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	client := Oauth2Application.Client(context.TODO(), tok)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	contents, _, err := utils.ReadAll(response.Body)

	if err != nil && err != io.EOF {
		ctx.JSON(500, err)
		return
	}

	parsingMap := make(map[string]string)

	err = json.Unmarshal([]byte(contents), &parsingMap)

	email, hasEmail := parsingMap["email"]

	if !hasEmail {
		ctx.JSON(500, "missing email!")
		return
	}

	account, err := db.GetAccountByEmail(&email)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(500, err)
		return
	}

	if account == nil {
		// 유저가 구글 로그인할 때 사용한 이메일을 다시 반환하여 프론트엔드에서 자동으로 이메일란을 채워놓을 수 있도록 함
		encrypted, err := utils.SignJWT(email, "email")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"email": encrypted})
	} else {
		token, err := utils.SignJWT(account.UserId, "user_id")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// OAuth로 로그인
		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}

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
	userpw := user.Password

	// Verify password
	if !utils.VerifyPassword(pw, userpw) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.SignJWT(user.UserId, "user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 이메일로 로그인
	c.JSON(http.StatusOK, gin.H{"token": token})
}
