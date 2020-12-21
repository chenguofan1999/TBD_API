package api

import (
	"fmt"
	"net/http"
	"tbd/model"
	"tbd/utils"

	"github.com/gin-gonic/gin"
)

func GetUserInfoByName(c *gin.Context) {
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

type bioInfo struct {
	Bio string `json:"bio" form:"bio"`
}

func UpdateUserBio(c *gin.Context) {
	// 得到登录用户 userID
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	userID := model.QueryUserIDWithName(loginUserName)

	if userID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not exist"})
		return
	}

	// 获取 JSON 中的参数
	var info bioInfo
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "bind error"})
		return
	}

	model.UpdateBio(userID, info.Bio)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpdateUserAvatar(c *gin.Context) {
	// 得到登录用户 userID
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)
	userID := model.QueryUserIDWithName(loginUserName)

	if userID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not exist"})
		return
	}

	// 读取文件
	imageFile, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expected Form-data: avatar"})
	}

	// 生成文件存储路径
	filePath := fmt.Sprintf("static/avatars/%s", utils.GenerateRandomFileName(imageFile.Filename))

	// 保存
	if err = c.SaveUploadedFile(imageFile, filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "upload error"})
		return
	}

	// 更新 UserInfo
	model.UpdateAvatar(userID, filePath)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
