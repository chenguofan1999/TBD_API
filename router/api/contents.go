package api

import (
	"fmt"
	"net/http"
	"strconv"
	"tbd/model"
	"tbd/utils"

	"github.com/gin-gonic/gin"
)

// GetContentsByQuerys 获取contents, 可选 query 项目:
// type : 查询类型：user / public / mine / following / like
// username : 查询用户内容、like时有效
// num : 公共内容时有效，默认 20 条
func GetContentsByQuerys(c *gin.Context) {
	queryType := c.Query("type")
	switch queryType {
	case "user":
		GetContentsByUserName(c)
	case "public":
		GetPublicContentsByNum(c)
	case "mine":
		GetMyContents(c)
	case "following":
		GetContentsOfMyFollowingUsers(c)
	case "like":
		GetLikedContentsByName(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BadRequest",
		})
		return
	}

}

// GetContentsByUserName 获取某用户的全部内容
func GetContentsByUserName(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BadRequest",
		})
		return
	}

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

// GetPublicContentsByNum 获取公共内容，query中用num表示条数
func GetPublicContentsByNum(c *gin.Context) {
	var number int
	var err error
	numStr := c.Query("num")
	if numStr == "" {
		number = 20
	} else {
		number, err = strconv.Atoi(numStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "BadRequest",
			})
			return
		}
	}

	contents, err := model.QueryPublicContents(number)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "BadRequest",
			})
			return
		}
	}
	c.JSON(http.StatusOK, contents)
}

// GetContentByContentID 根据 contentID 获取 content
func GetContentByContentID(c *gin.Context) {
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BadRequest",
		})
		return
	}

	content := model.QueryContentWithContentID(contentID)
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

// GetContentsOfMyFollowingUsers 获取当前用户关注者发布的内容
func GetContentsOfMyFollowingUsers(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	loginUserID := model.QueryUserIDWithName(loginUserName)

	contents, err := model.GetContentsOfFollowingUsersWithUserID(loginUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BadRequest",
		})
		return
	}

	c.JSON(http.StatusOK, contents)
}

// GetMyContents 获取当前用户发布的内容
func GetMyContents(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	contents, err := model.QueryContentsWithName(loginUserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BadRequest",
		})
		return
	}

	c.JSON(http.StatusOK, contents)
}

type TextContent struct {
	Title string `json:"title" form:"title"`
	Text  string `json:"text" form:"text"`
}

// PostContent 发布内容
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
		c.JSON(http.StatusNotFound, gin.H{"status": "Content not exist"})
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
