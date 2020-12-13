package api

import (
	//"fmt"
 	"net/http"
	"tbd/model"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

type LoginInfo struct {
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
}

func Login(c *gin.Context) {
	var info LoginInfo 
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"token": "",
			"status": "failed",
		})
		return ;
	}
	
	expected := model.QueryPasswordByName(info.Username)
	if(expected != "" && info.Password == expected) {
		tokenString := CreateTokenByName(info.Username)
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"status": "succeed",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": "", 
			"status": "failed",
		})	
	}
}

// used to create and parse token
var serversecret = []byte("randomkey")

// return "" if parse token failed
func GetNameByToken(tokenString string) (string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
  		return serversecret, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if username, ok := claims["username"].(string); ok {
				return username
			} 
		}
	}
	return ""
}

// return "" if create failed
func CreateTokenByName(name string) (string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": name,
	})
	tokenString, err := token.SignedString(serversecret)
	if err == nil {
		return tokenString
	} else {
		return ""
	}
}