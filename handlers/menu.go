package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
)

// GetCafeteriaMenus handles the GET /cafeteria endpoint
func GetCafeteriaMenus(c *gin.Context) {
	// Parse date from request query string
	dateParam := c.Query("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Get cafeteria menus from database
	menus, err := db.GetMenu(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cafeteria menus from database"})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// CreateCafeteriaMenu handles the POST /cafeteria endpoint
func CreateCafeteriaMenu(c *gin.Context) {
	// Parse request body
	var menu models.CafeteriaMenu
	err := c.BindJSON(&menu)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create menu in database
	_, err = db.CreateMenu(&menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cafeteria menu"})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

// UpdateCafeteriaMenu handles the PUT /cafeteria/:id endpoint
func UpdateCafeteriaMenu(c *gin.Context) {
	// Parse menu ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	// Get existing menu from database
	menu, err := db.GetMenuByID(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	// Parse request body
	var updatedMenu models.CafeteriaMenu
	err = c.BindJSON(&updatedMenu)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update existing menu with new data
	menu.Date = updatedMenu.Date
	menu.MealName = updatedMenu.MealName
	menu.Contents = updatedMenu.Contents

	// Update menu in database
	err = db.UpdateMenu(menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cafeteria menu"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteCafeteriaMenu handles the DELETE /cafeteria/:id endpoint
func DeleteCafeteriaMenu(c *gin.Context) {
	// Parse menu ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	// Delete menu from database
	err = db.DeleteMenu(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cafeteria menu"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
