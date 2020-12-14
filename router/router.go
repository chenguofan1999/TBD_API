package router

import (
	"tbd/router/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/user", api.GetLoginUser)
	router.GET("/user/:username", api.GetUserByName)
	router.GET("/users/:username/followers", api.GetFollowersByName)
	router.GET("/users/:username/following", api.GetFollowingByName)

	router.POST("/login", api.Login)

	router.GET("/contents/:contentID", api.GetContentByContentID)
	router.GET("/contents", api.GetContentsByName)

	router.GET("/comments/:commentID", api.GetCommentByCommentID)
	router.GET("/comments", api.GetCommentsByContentIDandFilter)

	router.GET("/replies/:commentID", api.GetRepliesByCommentID)

	// router.POST("/user", Store)
	// router.PUT("/user/:id", Update)
	// router.DELETE("/user/:id", Destroy)

	return router
}
