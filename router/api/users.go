package api

import (
	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

func GetUserByName(c *gin.Context) {
	username := c.Param("username")
	user := model.QueryUserWithName(username)

	// 没有找到,返回 404 not found
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetLoginUser(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	username := GetNameByToken(tokenString)
	if user := model.QueryUserWithName(username); user != nil {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}
}

/*
type LoginInfo struct {
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
}
*/
func CreateNewUser(c *gin.Context) {
	var info LoginInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "bind error"})
		return
	}

	if err := model.InsertUser(info.Username, info.Password); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
