package api

import (
	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

func GetContentsByName(c *gin.Context) {
	username := c.Param("username")

	// 确定此用户存在
	user := model.QueryUserWithName(username)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	contents := model.QueryContentsWithName(username)
	c.JSON(http.StatusOK, contents)
}
