package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/handlers"
	"github.com/username/schoolapp/middlewares"
	"github.com/username/schoolapp/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Connect()

	// Create new Gin router
	r := gin.Default()

	// Routes for handling authentication
	auth := r.Group("/auth")
	{
		oauth2Application := &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_ID"),
			ClientSecret: os.Getenv("OAUTH_SECRET"),
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:8080/auth/oauthsuccess",
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
			},
		}

		auth.GET("/oauthsuccess", func(ctx *gin.Context) {
			code := ctx.Query("code")

			tok, err := oauth2Application.Exchange(context.TODO(), code)
			if err != nil {
				ctx.JSON(500, err)
				return
			}

			client := oauth2Application.Client(context.TODO(), tok)

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
			if err != nil {
				ctx.JSON(500, err)
				return
			}

			SecretKey := os.Getenv("SECRET_KEY")
			if account == nil {
				// 유저가 구글 로그인할 때 사용한 이메일을 다시 반환하여 프론트엔드에서 자동으로 이메일란을 채워놓을 수 있도록 함
				encrypted, err := utils.EncryptJWT(email, "email", SecretKey)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"email": encrypted})
			} else {
				token, err := utils.EncryptJWT(account.UserId, "user_id", SecretKey)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"token": token})
			}
		})
		auth.PUT("/register", handlers.CreateAccount)
		auth.POST("/login", handlers.Login)
	}

	// Use authentication middleware for protected endpoints
	r.Use(middlewares.CheckAuthHeader)
	r.Use(middlewares.VerifyToken)

	// Routes for handling students
	students := r.Group("/students")
	{
		// Routes for handling timetables
		timetable := students.Group("/timetable")
		{
			timetable.GET("/lock", handlers.LockTimetable)
			timetable.GET("/unlock", handlers.UnLockTimetable)
			timetable.GET("", handlers.GetTimetableEntry)
			timetable.POST("", handlers.CreateTimetable)
			timetable.PUT("/:id", handlers.UpdateTimetable)
			timetable.DELETE("/:id", handlers.DeleteTimetable)
		}
	}

	admins := r.Group("/admins")
	{
		admins.GET("", handlers.GetAccountById)
		admins.PUT("/config", handlers.UpdateAccount)
	}

	// Routes for handling cafeteria menus
	cafeteriaMenus := r.Group("/cafeteria_menus")
	{
		cafeteriaMenus.GET("", handlers.GetCafeteriaMenus)
		cafeteriaMenus.POST("", handlers.CreateCafeteriaMenu)
		cafeteriaMenus.PUT("/:id", handlers.UpdateCafeteriaMenu)
		cafeteriaMenus.DELETE("/:id", handlers.DeleteCafeteriaMenu)
	}

	// Routes for handling checklists
	checklist := r.Group("/checklist")
	{
		checklist.GET("/lock", handlers.LockChecklist)
		checklist.GET("/unlock", handlers.UnLockChecklist)
		checklist.GET("", handlers.GetChecklist)
		checklist.POST("", handlers.CreateChecklist)
		checklist.PUT("/:id", handlers.UpdateChecklist)
		checklist.DELETE("/:id", handlers.DeleteChecklistItem)
	}

	events := r.Group("/events")
	{
		events.GET("/:months", handlers.GetEventsOfOneMonth)
	}

	r.GET("/map", handlers.GetMap)
	r.PUT("/map", handlers.PutMap)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Error running server")
	}
	db.close()
}
