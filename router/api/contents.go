package api

import (
	"net/http"
	"strconv"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

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
