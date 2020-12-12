package api

import (
	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

func GetUserByName(c *gin.Context) {
	username := c.Param("username")

	user := model.QueryWithName(username)

	// 没有找到,返回 404 not found
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
