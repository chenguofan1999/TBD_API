package router

import (
	"tbd/router/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/users/:username", api.GetUserByName)
	// router.POST("/user", Store)
	// router.PUT("/user/:id", Update)
	// router.DELETE("/user/:id", Destroy)

	return router
}
