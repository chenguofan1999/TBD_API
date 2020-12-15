package api

import (
	"fmt"
	"net/http"
	"strconv"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

// GetContentsByName 获取指定用户发表的所有content
func GetContentsByName(c *gin.Context) {
	username := c.Query("username")

	// 确定此用户存在
	user := model.QueryUserWithName(username)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No such user",
		})
		return
	}

	contents := model.QueryContentsWithName(username)
	c.JSON(http.StatusOK, contents)
}

// GetContentByContentID 根据 contentID 获取 content
func GetContentByContentID(c *gin.Context) {
	contentID64, err := strconv.ParseInt(c.Param("contentID"), 10, 32)
	if err != nil {
		panic(err)
	}
	contentID := int(contentID64)

	content := model.QueryContentWithContentID(contentID)
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

type TextContent struct {
	Title string `json:"title" form:"title"`
	Text  string `json:"text" form:"text"`
}

func PostContent(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	title := c.PostForm("title")
	text := c.PostForm("text")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expected Form-data"})
	}

	imageNum := model.QueryMaxImageID()
	imageURLs := make([]string, 0)
	imageFiles := form.File["imageFiles"]

	for _, file := range imageFiles {
		imageNum++
		filePath := fmt.Sprintf("static/images/%d-%s", imageNum, file.Filename)
		imageURLs = append(imageURLs, filePath)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "upload error"})
			return
		}
	}

	contentID, err := model.InsertContent(loginUserName, title, text, imageURLs)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"newContent": model.QueryContentWithContentID(contentID),
	})
}
