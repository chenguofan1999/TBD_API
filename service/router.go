package service

import (
	"net/http"
	"tbd/service/api"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(Cors())

	router.StaticFS("/static", http.Dir("./static"))

	router.GET("/users/:username", api.GetUserInfoByName)
	router.GET("/users/:username/followers", api.GetFollowersByName)
	router.GET("/users/:username/following", api.GetFollowingByName)
	router.GET("/users/:username/likes", api.GetLikedContentsByName)

	router.POST("/users/login", api.Login)
	router.POST("/users/signup", api.CreateNewUser)

	router.GET("/contents/:contentID", api.GetContentByContentID)
	router.GET("/contents", api.GetContentsByQuerys)

	router.GET("/comments/:commentID", api.GetCommentByCommentID)
	router.GET("/comments", api.GetCommentsByContentIDandFilter)

	router.GET("/comments/:commentID/replies", api.GetRepliesByCommentID)

	// Following API based on current login user, please include a token in request header
	router.PUT("/user/info/bio", api.UpdateUserBio)
	router.PUT("/user/info/avatar", api.UpdateUserAvatar)
	router.GET("/user", api.GetLoginUser)
	router.PUT("/user/following/:username", api.FollowUser)
	router.DELETE("/user/following/:username", api.UnfollowUser)
	router.GET("/user/following/:username", api.GetFollowStateByUsername)

	router.POST("/contents", api.PostContent)
	router.DELETE("/contents/:contentID", api.DeleteContent)

	router.POST("/comments", api.PostComment)
	router.POST("/comments/:commentID/replies", api.PostReply) // todo:
	router.DELETE("/comments/:commentID", api.DeleteComment)

	router.PUT("/contents/:contentID/like", api.LikeContentByContentID)
	router.DELETE("/contents/:contentID/like", api.CancelLikeContentByContentID)
	router.GET("/contents/:contentID/like", api.CheckLiked)

	return router
}
