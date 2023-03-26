package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
)

// Get all checklists for a student
func GetChecklists(c *gin.Context) {
	studentID := c.Param("student_id")

	checklists, err := db.GetChecklists(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, checklists)
}

// Create a new checklist
func CreateChecklist(c *gin.Context) {
	var checklist models.Checklist

	if err := c.ShouldBindJSON(&checklist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateChecklist(&checklist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.CreateChecklist(&checklist)
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

// Update a checklist
func UpdateChecklist(c *gin.Context) {
	var checklist models.Checklist

	if err := c.ShouldBindJSON(&checklist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateChecklist(&checklist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateChecklist(&checklist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Delete a checklist
func DeleteChecklist(c *gin.Context) {
	// Parse checklist ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete checklist from database
	err = db.DeleteChecklist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
