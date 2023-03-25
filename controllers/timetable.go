package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
)

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
	id := c.Param("id")

	err := db.DeleteTimetable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
