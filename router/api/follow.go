package api

import (
	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

// GetFollowersByName : 得到指定用户的关注者
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

// GetFollowingByName : 得到指定用户关注的人
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

// FollowUser : 当前登录用户关注目标用户
func FollowUser(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	if loginUser := model.QueryUserWithName(loginUserName); loginUser == nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	if err := model.InsertFollowRelationByName(loginUserName, c.Param("username")); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// UnfollowUser : 当前登录用户取消关注目标用户
func UnfollowUser(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	if loginUser := model.QueryUserWithName(loginUserName); loginUser == nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	if err := model.DeleteFollowRelationByName(loginUserName, c.Param("username")); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
