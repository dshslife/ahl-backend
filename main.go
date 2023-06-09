package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3"
	"github.com/joho/godotenv"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/handlers"
	"github.com/username/schoolapp/middlewares"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect()
	models.InitAllergies()
	utils.InitKeys()

	handlers.Oauth2Application = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_ID"),
		ClientSecret: os.Getenv("OAUTH_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/oauthsuccess",
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}

	// Create new Gin router
	r := gin.Default()

	// 모든 JWT 토큰은 클라이언트 공개 키로 서명을 확인해야 함
	// 거의 모든 엔드포인트가 JWT 토큰을 사용하니, 거의 다 클라이언트 공개 키를 알아야 함
	// 그래서 middleware로 이를 강제하는게 적당함
	r.Use(middlewares.EnforceClientJWTPubKey)

	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	loginData.Email = "iusedlinux@gmail.com"
	loginData.Password = "drowssap"

	contents, err := os.ReadFile("./user_keys/private.pem")
	if err != nil {
		panic(err.Error())
	}

	block, _ := pem.Decode(contents)

	userPrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err.Error())
	}

	signed, err := utils.SignJWTWithKey(loginData, "account", userPrivateKey.(*rsa.PrivateKey))
	if err != nil {
		panic(err.Error())
	}

	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: &utils.PRIVATE.PublicKey}, nil)
	if err != nil {
		panic(err.Error())
	}

	obj, err := encrypter.Encrypt([]byte(signed))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(obj.FullSerialize())

	// 아래 미들 웨어들은 계정 로그인 여부를 검사 및 강제하는데
	// /auth 엔드포인트들은 아직 로그인을 안했을 수도 있기 때문에 무시함
	// /publickey는 언제나 공개이기에 무시함
	r.Use(middlewares.CheckAuthHeader)
	r.Use(middlewares.VerifyToken)

	r.GET("/publickey", handlers.GetPublicKey)

	// Routes for handling authentication
	auth := r.Group("/auth")
	{
		auth.GET("/oauth", handlers.OAuthRedirect)
		auth.GET("/oauthsuccess", handlers.OnOAuth)
		// TODO CreateAccount로 다시 바꾸기
		auth.PUT("/register", handlers.CreateAccountUnsafe)
		auth.POST("/login", handlers.Login)
	}

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
	db.Close()
}
