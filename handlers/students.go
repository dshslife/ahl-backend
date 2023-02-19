package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
)

// GetStudents handles the GET /students endpoint
func GetStudents(c *gin.Context) {
	// Get students from database
	students, err := db.GetStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students from database"})
		return
	}

	c.JSON(http.StatusOK, students)
}
