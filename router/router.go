package router

import (
	"net/http"
	"tbd/router/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.StaticFS("/static", http.Dir("./static"))

	router.GET("/users/:username", api.GetUserByName)
	router.GET("/users/:username/followers", api.GetFollowersByName)
	router.GET("/users/:username/following", api.GetFollowingByName)

	router.POST("/login", api.Login)
	router.POST("/signup", api.CreateNewUser)

	router.GET("/contents/:contentID", api.GetContentByContentID)
	router.GET("/contents", api.GetContentsByName)

	router.GET("/comments/:commentID", api.GetCommentByCommentID)
	router.GET("/comments", api.GetCommentsByContentIDandFilter)

	router.GET("/replies/:commentID", api.GetRepliesByCommentID)

	// Following API based on current login user, please include a token in request header
	router.GET("/user", api.GetLoginUser)
	router.PUT("/user/following/:username", api.FollowUser)
	router.DELETE("/user/following/:username", api.UnfollowUser)

	router.POST("/contents", api.PostContent)
	router.POST("/comments", api.PostComment)
	router.POST("/replies", api.PostReply)

	// router.POST("/user", Store)
	// router.PUT("/user/:id", Update)
	// router.DELETE("/user/:id", Destroy)

	return router
}
