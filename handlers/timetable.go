package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"net/http"
	"strconv"
)

// LockTimetable locks the timetable
func LockTimetable(c *gin.Context) {
	fetched, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	studentId := fetched.(uuid.UUID)
	user, err := db.GetAccountById(&studentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.PermissionInfo.GetLevel() != models.STUDENT {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested user is not student"})
	}
	info := user.PermissionInfo.(models.StudentInfo)
	info.Timetable.IsPublic = false

	c.JSON(http.StatusOK, info.Timetable)
}

// UnLockTimetable unlocks the timetable
func UnLockTimetable(c *gin.Context) {
	fetched, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	studentId := fetched.(uuid.UUID)
	user, err := db.GetAccountById(&studentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.PermissionInfo.GetLevel() != models.STUDENT {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested user is not student"})
	}
	info := user.PermissionInfo.(models.StudentInfo)
	info.Timetable.IsPublic = true

	c.JSON(http.StatusOK, info.Timetable)
}

// GetTimetableEntry handles the GET /timetable endpoint
func GetTimetableEntry(c *gin.Context) {
	// Get user ID from request context
	fetched, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Get timetable from database
	dbId := models.DbId(fetched.(int))
	timetable, err := db.GetTimeTableEntry(dbId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timetable from database"})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

// CreateTimetable handles the POST /timetable endpoint
func CreateTimetable(c *gin.Context) {
	// Get user ID from request context
	fetched, exists := c.Get("teacherID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Parse request body
	var lesson models.TimetableEntry
	err := c.BindJSON(&lesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set user ID for new lesson
	lesson.TeacherId = fetched.(uuid.UUID)

	// Create lesson in database
	_, err = db.CreateTimetable(&lesson)
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
	_, err = db.GetTimeTableEntry(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	// Parse request body
	var updatedLesson models.TimetableEntry
	err = c.BindJSON(&updatedLesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update lesson in database
	err = db.UpdateTimetable(&updatedLesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lesson"})
		return
	}

	c.JSON(http.StatusOK, updatedLesson)
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
	err = db.DeleteTimetable(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
