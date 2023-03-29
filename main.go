package main

import (
	"fmt"
	"github.com/username/schoolapp/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		students.GET("", handlers.GetStudents)
		students.PUT("/config", handlers.UpdateStudent)
	}

	teachers := r.Group("/teachers")
	{
		teachers.GET("", handlers.GetTeachers)
		teachers.PUT("/config", handlers.UpdateTeacher)
	}

	admins := r.Group("/admins")
	{
		admins.GET("/admin_only", handlers.GetAdmins)
		admins.PUT("/admin_only/config", handlers.UpdateAdmin)
	}

	// Routes for handling timetables
	timetable := r.Group("/timetable")
	{
		timetable.GET("/lock", handlers.LockTimetable)
		timetable.GET("/unlock", handlers.UnLockTimetable)
		timetable.GET("", handlers.GetTimetable)
		timetable.POST("", handlers.CreateTimetable)
		timetable.PUT("/:id", handlers.UpdateTimetable)
		timetable.DELETE("/:id", handlers.DeleteTimetable)
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
		checklist.POST("", handlers.CreateChecklistItem)
		checklist.PUT("/:id", handlers.UpdateChecklistItem)
		checklist.DELETE("/:id", handlers.DeleteChecklistItem)
	}

	events := r.Group("/events")
	{
		events.GET("/:months", handlers.GetEvents)
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
