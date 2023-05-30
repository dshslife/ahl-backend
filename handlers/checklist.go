package handlers

import (
	"github.com/google/uuid"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
)

// LockChecklist locks the checklist
func LockChecklist(c *gin.Context) {
	studentID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	temp := studentID.(uuid.UUID)
	checklist, err := db.GetChecklistsOfStudent(&temp)
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
	studentID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	temp := studentID.(uuid.UUID)
	checklist, err := db.GetChecklistsOfStudent(&temp)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}

	// Get checklist items from database
	temp := userID.(uuid.UUID)
	items, err := db.GetChecklistsOfStudent(&temp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get checklist items from database"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateChecklist handles the POST /checklist endpoint
func CreateChecklist(c *gin.Context) {
	// Get user ID from request context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
		return
	}
	temp := userID.(uuid.UUID)

	// Parse request body
	var checklist models.Checklist
	err := c.BindJSON(&checklist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	existing, _ := db.GetChecklistsOfStudent(&temp)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Checklist already exists"})
		return
	}

	// Execute query to insert checklist items
	for _, items := range checklist.Items {
		items.Complete = false
		checklist.Items = append(checklist.Items, items)
	}

	// Set user ID and completed status for new checklist
	checklist.StudentId = userID.(uuid.UUID)
	checklist.Title = "checklist"

	// Create checklist checklist in database
	_, err = db.CreateChecklist(&checklist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checklist checklist"})
		return
	}

	c.JSON(http.StatusCreated, checklist)
}

// UpdateChecklist handles the PUT /checklist/:id endpoint
func UpdateChecklist(c *gin.Context) {
	// Parse toUpdate ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toUpdate ID"})
		return
	}

	// Get existing toUpdate from database
	toUpdate, err := db.GetChecklistsById(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MenuEntry not found"})
		return
	}

	// Parse request body
	var updatedItem models.Checklist
	err = c.BindJSON(&updatedItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update existing toUpdate with new data
	toUpdate.Title = updatedItem.Title
	toUpdate.Items = updatedItem.Items

	// Update toUpdate in database
	err = db.UpdateChecklist(toUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update checklist toUpdate"})
		return
	}

	c.JSON(http.StatusOK, toUpdate)
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
	err = db.DeleteChecklist(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete checklist item"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
