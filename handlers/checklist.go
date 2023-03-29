package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
)

// LockChecklist locks the checklist
func LockChecklist(c *gin.Context) {
	studentID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	checklist, err := db.GetChecklists(studentID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Execute query to insert checklist items
	for _, item := range checklist.Items {
		item.IsPublic = false
		checklist.Items = append(checklist.Items, item)
	}

	c.JSON(http.StatusOK, checklist)
}

// UnLockChecklist unlocks the checklist
func UnLockChecklist(c *gin.Context) {
	studentID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	checklist, err := db.GetChecklists(studentID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Execute query to insert checklist items
	for _, item := range checklist.Items {
		item.IsPublic = true
		checklist.Items = append(checklist.Items, item)
	}

	c.JSON(http.StatusOK, checklist)
}

// GetChecklist handles the GET /checklist endpoint
func GetChecklist(c *gin.Context) {
	// Get user ID from request context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Get checklist items from database
	items, err := db.GetChecklistItems(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get checklist items from database"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateChecklistItem handles the POST /checklist endpoint
func CreateChecklistItem(c *gin.Context) {
	// Get user ID from request context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Parse request body
	var item models.Checklist
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Execute query to insert checklist items
	for _, items := range item.Items {
		items.Complete = false
		item.Items = append(item.Items, items)
	}

	// Set user ID and completed status for new item
	item.UserID = userID.(string)

	// Create checklist item in database
	err = db.CreateChecklistItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checklist item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateChecklistItem handles the PUT /checklist/:id endpoint
func UpdateChecklistItem(c *gin.Context) {
	// Parse item ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get existing item from database
	item, err := db.GetChecklistItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Parse request body
	var updatedItem models.Checklist
	err = c.BindJSON(&updatedItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update existing item with new data
	item.Title = updatedItem.Title
	item.Items = updatedItem.Items

	// Update item in database
	err = db.UpdateChecklistItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update checklist item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteChecklistItem handles the DELETE /checklist/:id endpoint
func DeleteChecklistItem(c *gin.Context) {
	// Parse item ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Delete item from database
	err = db.DeleteChecklistItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete checklist item"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
