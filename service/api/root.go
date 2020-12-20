package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBriefAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"current_user_url": "http://159.75.1.231:5005/user",
		"user_url":         "http://159.75.1.231:5005/users/{username}",
		"following_url":    "http://159.75.1.231:5005/users/{username}/followint",
		"followers_url":    "http://159.75.1.231:5005/users/{username}/followers",
		"content_url":      "http://159.75.1.231:5005/content/{contentID}",
		"contents_url":     "http://159.75.1.231:5005/contents",
		"comment_url":      "http://159.75.1.231:5005/comments/{commentID}",
		"comments_url":     "http://159.75.1.231:5005/comments",
	})
}
