package api

import (
	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

func GetFollowersByName(c *gin.Context) {
	username := c.Param("username")

	// 确定此用户存在
	user := model.QueryUserWithName(username)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	followers := model.QueryFollowersWithName(username)
	c.JSON(http.StatusOK, followers)
}

func GetFollowingByName(c *gin.Context) {
	username := c.Param("username")

	// 确定此用户存在
	user := model.QueryUserWithName(username)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	followers := model.QueryFollowingWithName(username)
	c.JSON(http.StatusOK, followers)
}
