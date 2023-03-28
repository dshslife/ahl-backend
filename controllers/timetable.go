package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
)

// Lock the timetable
func LockTimetable(c *gin.Context) {
	studentID := c.Param("student_id")

	var timetable models.Timetable
	timetables, err := db.GetTimetable(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	timetable.IsPublic = false
	timetables = append(timetables, timetable)

	c.JSON(http.StatusOK, timetables)
}

// UnLock the timetable
func UnLockTimetable(c *gin.Context) {
	studentID := c.Param("student_id")

	var timetable models.Timetable
	timetables, err := db.GetTimetable(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	timetable.IsPublic = true
	timetables = append(timetables, timetable)

	c.JSON(http.StatusOK, timetables)
}

// Get all timetables for a student
func GetTimetable(c *gin.Context) {
	studentID := c.Param("student_id")

	timetables, err := db.GetTimetable(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, timetables)
}

// Create a new timetable
func CreateTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := c.ShouldBindJSON(&timetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateTimetable(&timetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.CreateTimetable(&timetable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// Update a timetable
func UpdateTimetable(c *gin.Context) {
	var timetable models.Timetable

	if err := c.ShouldBindJSON(&timetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateTimetable(&timetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateTimetable(&timetable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Delete a timetable
func DeleteTimetable(c *gin.Context) {
	// Parse lesson ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete lesson from database
	err = db.DeleteTimetable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
