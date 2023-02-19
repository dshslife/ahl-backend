package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
)

// GetTimetable handles the GET /timetable endpoint
func GetTimetable(c *gin.Context) {
	// Get user ID from request context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Get timetable from database
	timetable, err := db.GetTimetable(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timetable from database"})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

// CreateTimetable handles the POST /timetable endpoint
func CreateTimetable(c *gin.Context) {
	// Get user ID from request context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Parse request body
	var lesson models.Lesson
	err := c.BindJSON(&lesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set user ID for new lesson
	lesson.UserID = userID.(int)

	// Create lesson in database
	err = db.CreateLesson(&lesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lesson"})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// UpdateTimetable handles the PUT /timetable/:id endpoint
func UpdateTimetable(c *gin.Context) {
	// Parse lesson ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	// Get existing lesson from database
	lesson, err := db.GetLessonByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	// Parse request body
	var updatedLesson models.Lesson
	err = c.BindJSON(&updatedLesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update existing lesson with new data
	lesson.Name = updatedLesson.Name
	lesson.Teacher = updatedLesson.Teacher
	lesson.Location = updatedLesson.Location
	lesson.StartTime = updatedLesson.StartTime
	lesson.EndTime = updatedLesson.EndTime
	lesson.Day = updatedLesson.Day

	// Update lesson in database
	err = db.UpdateLesson(lesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lesson"})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// DeleteTimetable handles the DELETE /timetable/:id endpoint
func DeleteTimetable(c *gin.Context) {
	// Parse lesson ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	// Delete lesson from database
	err = db.DeleteLesson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
