package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
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

// GetAccountById 리퀘스트 바디로 UUID String이 주어지면 이를 분석하고 이에 맞는 계정 반환
// 근데 리퀘스트부터 UUID String이 아니라 바이트 배열로 보내면 분석을 건너뛰어도 되는데 이게 더 나을려나?
func GetAccountById(c *gin.Context) {
	uuidString, _, err := utils.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	id, err := uuid.Parse(uuidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

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
	buf := new(strings.Builder)
	_, err := io.Copy(buf, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contents := buf.String()
	decrypted, err := utils.DecryptJWT(&contents, "account")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var account models.Account
	err = json.Unmarshal([]byte(decrypted), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidateNewAccount(&account); err != nil {
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
