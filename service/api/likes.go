package api

import (
	"net/http"
	"strconv"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

func GetLikedContentsByName(c *gin.Context) {
	username := c.Param("username")
	userID := model.QueryUserIDWithName(username)

	contents, err := model.QueryLikedContentsWithUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, contents)
}

func LikeContentByContentID(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	loginUserID := model.QueryUserIDWithName(loginUserName)

	// 从 URL 中得到 contentID
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = model.InsertLikeRelation(loginUserID, contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func CancelLikeContentByContentID(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	loginUserID := model.QueryUserIDWithName(loginUserName)

	// 从 URL 中得到 contentID
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 执行删除
	err = model.DeleteLikeRelation(loginUserID, contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func CheckLiked(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	loginUserID := model.QueryUserIDWithName(loginUserName)

	// 从 URL 中得到 contentID
	contentID, err := strconv.Atoi(c.Param("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	liked, err := model.QueryHasLiked(loginUserID, contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
	})
}
