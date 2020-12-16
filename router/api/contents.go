package api

import (
	"fmt"
	"net/http"
	"strconv"
	"tbd/model"
	"tbd/utils"

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

	contents, err := model.QueryContentsWithName(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

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

	imageURLs := make([]string, 0)
	imageFiles := form.File["images"]

	for _, file := range imageFiles {
		filePath := fmt.Sprintf("static/images/%s", utils.GenerateRandomFileName(file.Filename))
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

// DeleteContent : 删除内容，在 URL 的 Path 中写入 contentID
func DeleteContent(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	// 获取 URL 中 contentID 并验证格式是否正确
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request"})
		return
	}

	// 验证 content 是否存在，存在即获取
	content := model.QueryContentWithContentID(contentID)
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Bad Request"})
		return
	}

	// 验证是否为该用户所发
	if content.Author.Username != loginUserName {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	// 执行删除
	err = model.DeleteContentWithContentID(contentID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
