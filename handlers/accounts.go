package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"net/http"
	"strconv"
)

func GetAllAccounts(c *gin.Context) {
	accounts, err := db.GetAllAccounts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func GetAccountById(c *gin.Context) {
	id := models.UserId(c.Param("id"))

	account, err := db.GetAccountById(&id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "account not found",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

func CreateAccount(c *gin.Context) {
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateAccount(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.CreateAccount(&account)
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

func UpdateAccount(c *gin.Context) {
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.ValidateAccount(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "Account updated")
}

func DeleteAccount(c *gin.Context) {
	// Parse account ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete account from database
	err = db.DeleteAccount(models.DbId(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "Account deleted")
}
