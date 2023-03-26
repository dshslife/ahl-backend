package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
)

// Get all cafeteria menus
func GetCafeteriaMenus(c *gin.Context) {
	menus, err := db.GetCafeteriaMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// Create a new cafeteria menu
func CreateCafeteriaMenu(c *gin.Context) {
	var menu models.CafeteriaMenu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateCafeteriaMenu(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.CreateCafeteriaMenu(&menu)
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

// Update a cafeteria menu
func UpdateCafeteriaMenu(c *gin.Context) {
	var menu models.CafeteriaMenu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateCafeteriaMenu(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateCafeteriaMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Delete a cafeteria menu
func DeleteCafeteriaMenu(c *gin.Context) {
	// Parse menu ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete menu from database
	err = db.DeleteCafeteriaMenu(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
