package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/username/schoolapp/db"
	"github.com/username/schoolapp/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetMap(c *gin.Context) {
	fetched, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "user id is not supplied to auth header")
		return
	}

	userId := models.UserId(fetched.(string))
	account, err := db.GetAccountById(&userId)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	var filename string
	switch account.PermissionInfo.GetLevel() {
	case models.STUDENT:
		info := account.PermissionInfo.(models.StudentInfo)
		filename = string(info.SchoolId)
	case models.TEACHER:
		info := account.PermissionInfo.(models.TeacherInfo)
		filename = string(info.SchoolId)
	default:
		c.String(http.StatusForbidden, "only student and teacher account can upload maps")
		return
	}

	found := false
	err = filepath.Walk("./maps", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !found && strings.HasPrefix(info.Name(), filename) {
			found = true
			if filepath.Ext(info.Name()) == ".png" {
				c.Header("Content-Type", "image/png")
			} else {
				c.Header("Content-Type", "image/jpeg")
			}
			http.ServeFile(c.Writer, c.Request, info.Name())
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, "map file not found")
		return
	}
}

func PutMap(c *gin.Context) {
	fetched, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "user id is not supplied to auth header")
		return
	}

	userId := models.UserId(fetched.(string))
	account, err := db.GetAccountById(&userId)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	var filename string
	switch account.PermissionInfo.GetLevel() {
	case models.STUDENT:
		info := account.PermissionInfo.(models.StudentInfo)
		filename = string(info.SchoolId)
	case models.TEACHER:
		info := account.PermissionInfo.(models.TeacherInfo)
		filename = string(info.SchoolId)
	default:
		c.String(http.StatusForbidden, "only student and teacher account can upload maps")
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if header.Size > 1024*1024*10 {
		c.String(http.StatusBadRequest, "Size too big!")
		return
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	extension := filepath.Ext(header.Filename)
	if extension != ".jpg" && extension != ".png" && extension != ".jpeg" {
		c.JSON(http.StatusBadRequest, "invalid file extension, only use .jpg, .png, or .jpeg")
	}
	path := "./maps/" + filename + extension

	err = os.WriteFile(path, contents, 777)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
}
