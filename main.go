package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/username/schoolapp/controllers"
	"github.com/username/schoolapp/middlewares"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create new Gin router
	r := gin.Default()

	// Use authentication middleware for protected endpoints
	authMiddleware := middlewares.NewAuthMiddleware()
	r.Use(authMiddleware.MiddlewareFunc())

	// Routes for handling students
	students := r.Group("/students")
	{
		students.GET("", controllers.GetStudents)
	}

	// Routes for handling timetables
	timetable := r.Group("/timetable")
	{
		timetable.GET("", controllers.GetTimetable)
		timetable.POST("", controllers.CreateTimetable)
		timetable.PUT("/:id", controllers.UpdateTimetable)
		timetable.DELETE("/:id", controllers.DeleteTimetable)
	}

	// Routes for handling cafeteria menus
	cafeteriaMenus := r.Group("/cafeteria_menus")
	{
		cafeteriaMenus.GET("", controllers.GetCafeteriaMenus)
		cafeteriaMenus.POST("", controllers.CreateCafeteriaMenu)
		cafeteriaMenus.PUT("/:id", controllers.UpdateCafeteriaMenu)
		cafeteriaMenus.DELETE("/:id", controllers.DeleteCafeteriaMenu)
	}

	// Routes for handling checklists
	checklist := r.Group("/checklist")
	{
		checklist.POST("", controllers.CreateChecklistItem)
		checklist.PUT("/:id", controllers.UpdateChecklistItem)
		checklist.DELETE("/:id", controllers.DeleteChecklistItem)
	}

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Error running server")
	}
}
