package handlers

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// CreateAccountUnsafe 테스트 전용, TODO 무조건 삭제할 것!
func CreateAccountUnsafe(c *gin.Context) {
	body, _, err := utils.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var account models.Account
	asMap := make(map[string]interface{})

	_ = json.Unmarshal([]byte(body), &asMap)

	necessary := [5]string{"name", "email", "password", "permission_level", "permission"}
	for _, value := range necessary {
		if asMap[value] == "" {
			c.String(http.StatusBadRequest, "%s required", value)
			return
		}
	}
	account.Name = asMap["name"].(string)
	account.Email = asMap["email"].(string)
	password := asMap["password"].(string)

	_, err = db.GetAccountByEmail(&account.Email)
	if err != sql.ErrNoRows {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.Status(http.StatusForbidden)
		}
		return
	}

	passwordHash, err := utils.HashPassword([]byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account.Password = passwordHash
	level := asMap["permission_level"]

	permissionInfoRaw := asMap["permission"]
	permissionInfo, ok := permissionInfoRaw.(map[string]string)
	if !ok {
		_, isEmptyMap := permissionInfoRaw.(interface{})
		if !isEmptyMap {
			c.String(http.StatusBadRequest, "invalid permission info")
			return
		}
		permissionInfo = make(map[string]string)
	}
	switch level {
	case "student":
		{
			permission := models.StudentInfo{}

			permission.Class, err = strconv.Atoi(permissionInfo["class"])
			if err != nil {
				c.String(http.StatusBadRequest, "class should be number")
				return
			}
			permission.Grade, err = strconv.Atoi(permissionInfo["grade"])
			if err != nil {
				c.String(http.StatusBadRequest, "grade should be number")
				return
			}
			permission.Number, err = strconv.Atoi(permissionInfo["number"])
			if err != nil {
				c.String(http.StatusBadRequest, "class number should be number, duh")
				return
			}
			permission.SchoolId = models.SchoolId(permissionInfo["school_id"])
			account.PermissionInfo = permission
		}
	case "teacher":
		{
			permission := models.TeacherInfo{}
			permission.SchoolId = models.SchoolId(permissionInfo["school_id"])
			account.PermissionInfo = permission
		}
	case "admin":
		{
			permission := models.AdminInfo{}
			account.PermissionInfo = permission
		}
	default:
		c.String(http.StatusBadRequest, "unknown permission level")
		return
	}

	if err := db.ValidateNewAccount(&account); err != nil {
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

func CreateAccount(c *gin.Context) {
	contents, _, err := utils.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	publicKey, _ := c.Get("client_key")

	decrypted, err := utils.ParseJWT(&contents, "account", publicKey.(*rsa.PublicKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var account models.Account
	err = json.Unmarshal([]byte(decrypted.(string)), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.ValidateNewAccount(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = db.CreateAccount(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
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
