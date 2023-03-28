package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"net/http"
	"strconv"
)

// Get all teachers
func GetTeachers(c *gin.Context) {
	teachers, err := db.GetAllTeachers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// Get a teacher by ID
func GetTeacher(c *gin.Context) {
	id := c.Param("id")

	teacher, err := db.GetTeacherByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "teacher not found",
		})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// Create a new teacher
func CreateTeacher(c *gin.Context) {
	var teacher models.Teacher

	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateTeacher(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.CreateTeacher(&teacher)
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

// Update a teacher
func Updateteacher(c *gin.Context) {
	var teacher models.Teacher

	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateTeacher(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateTeacher(&teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Delete a teacher
func DeleteTeacher(c *gin.Context) {
	// Parse teacher ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete teacher from database
	err = db.DeleteTeacher(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
