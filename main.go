package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/handlers"
	"github.com/username/schoolapp/middlewares"
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

	// Create new Gin router
	r := gin.Default()

	// Use authentication middleware for protected endpoints
	r.Use(middlewares.Auth())
	r.Use(middlewares.VerifyToken)

	// Routes for handling authentication
	auth := r.Group("/auth")
	{
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
	db.close()
}
